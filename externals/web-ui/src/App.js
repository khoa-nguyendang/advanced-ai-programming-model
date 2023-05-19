import MuiAlert from '@mui/material/Alert';
import Card from '@mui/material/Card';
import CardContent from '@mui/material/CardContent';
import Slide from '@mui/material/Slide';
import Snackbar from '@mui/material/Snackbar';
import Typography from '@mui/material/Typography';
import Grid from '@mui/material/Unstable_Grid2';
import * as faceDetection from '@tensorflow-models/face-detection';
import * as tfjsWasm from '@tensorflow/tfjs-backend-wasm';
import '@tensorflow/tfjs-backend-webgl';
import '@tensorflow/tfjs-backend-webgpu';
import { useRef, useState } from 'react';
import { useEffectOnce } from 'react-use';
import './App.css';
import { Camera } from './Camera';
import banner from './fit-logo-kem-truong.png';
import { setupDatGui } from './option_panel';
import { createDetector, STATE } from './shared/params';
import { setBackendAndEnvFlags } from './shared/util';
tfjsWasm.setWasmPaths(`https://cdn.jsdelivr.net/npm/@tensorflow/tfjs-backend-wasm@${tfjsWasm.version_wasm}/dist/`);



function App(props) {
    const [appState, setAppState] = useState({
        trackId: "",
        alertType: "",
        message: "",
        open: false,
        userInfo: {}
    });
    const video = useRef(null);
    useEffectOnce(() => {
        console.log("use effect one");
        let detector, camera, stats;
        let startInferenceTime, numInferences = 0;
        let inferenceTimeSum = 0, lastPanelUpdate = 0;
        let rafId;
        const model = faceDetection.SupportedModels.MediaPipeFaceDetector;
        const detectorConfig = {
            runtime: 'tfjs', // or 'tfjs' 'mediapipe'
        }

        const loadDetector = async () => {
            detector = await faceDetection.createDetector(model, detectorConfig);
        };

        loadDetector();

        const verifyUserInfo = async function (image) {
            image = image.replace("data:image/jpeg;base64,", "");
            console.log("calling the api verifyUserInfo");
            let body = {
                images: [{
                    data: image,
                    imageId: Date.now().toLocaleString()
                }]
            }
            await fetch('http://localhost:8081/api/v1/user/verify', {
                method: 'POST',
                body: JSON.stringify(body),
                headers: {
                    'Content-type': 'application/json; charset=UTF-8',
                },
            })
                .then((response) => response.json())
                .then((data) => {
                    setAppState((prevState) => { return { ...prevState, userInfo: data } });
                    if (data && data.code == 200) {
                        printed = 1;
                        setTimeout(function () {
                            printed = 0;
                        }, 5000);
                    } else {
                        setTimeout(function () {
                            printed = 0;
                        }, 500);
                    }
                })
                .catch((err) => {
                    console.log(err.message);
                });
        }

        const checkGuiUpdate = async function () {
            if (STATE.isTargetFPSChanged || STATE.isSizeOptionChanged) {
                camera = await Camera.setupCamera(STATE.camera);
                STATE.isTargetFPSChanged = false;
                STATE.isSizeOptionChanged = false;
            }

            if (STATE.isModelChanged || STATE.isFlagChanged || STATE.isBackendChanged) {
                STATE.isModelChanged = true;

                window.cancelAnimationFrame(rafId);

                if (detector != null) {
                    detector.dispose();
                }

                if (STATE.isFlagChanged || STATE.isBackendChanged) {
                    await setBackendAndEnvFlags(STATE.flags, STATE.backend);
                }

                try {
                    detector = await createDetector(STATE.model);
                } catch (error) {
                    detector = null;
                    console.log("createDetector issue: ", error);
                }

                STATE.isFlagChanged = false;
                STATE.isBackendChanged = false;
                STATE.isModelChanged = false;
            }
        }

        const beginEstimateFaceStats = function () {
            startInferenceTime = (performance || Date).now();
        }

        const endEstimateFaceStats = function () {
            const endInferenceTime = (performance || Date).now();
            inferenceTimeSum += endInferenceTime - startInferenceTime;
            ++numInferences;

            const panelUpdateMilliseconds = 1000;
            if (endInferenceTime - lastPanelUpdate >= panelUpdateMilliseconds) {
                // const averageInferenceTime = inferenceTimeSum / numInferences;
                inferenceTimeSum = 0;
                numInferences = 0;
                // stats.customFpsPanel.update(1000.0 / averageInferenceTime, 120 /* maxValue */);
                lastPanelUpdate = endInferenceTime;
            }
        }

        let printed = 0;
        const renderResult = async function () {
            if (camera.video.readyState < 2) {
                await camera.waitingRenderVideoSource();
            }

            let faces = null;

            // Detector can be null if initialization failed (for example when loading
            // from a URL that does not exist).
            if (detector != null) {
                // FPS only counts the time it takes to finish estimateFaces.
                beginEstimateFaceStats();

                // Detectors can throw errors, for example when using custom URLs that
                // contain a model that doesn't provide the expected output.
                try {
                    faces = await detector.estimateFaces(camera.video, { flipHorizontal: false });
                } catch (error) {
                    detector.dispose();
                    detector = null;
                    console.log(" render Result issue: ", error);
                }

                endEstimateFaceStats();
            }

            camera.drawCtx();

            // The null check makes sure the UI is not in the middle of changing to a
            // different model. If during model change, the result is from an old model,
            // which shouldn't be rendered.
            if (faces && faces.length > 0 && !STATE.isModelChanged) {
                camera.drawResults(faces, STATE.modelConfig.boundingBox, false);
                //send picture to server for verify
                if (!printed) {
                    let image = camera.takeShoot();
                    await verifyUserInfo(image);
                    // console.log(faces);
                    printed++;

                }
            }
        }
        const aiapp = async function () {
            // Gui content will change depending on which model is in the query string.
            const urlParams = new URLSearchParams("?model=mediapipe_face_detector");
            if (!urlParams.has('model')) {
                alert('Cannot find model in the query string.');
                return;
            }

            await setupDatGui(urlParams);

            // stats = setupStats();

            camera = await Camera.setupCamera(STATE.camera);

            await setBackendAndEnvFlags(STATE.flags, STATE.backend);

            detector = await createDetector();

            renderPrediction();
        };
        const renderPrediction = async function () {
            await checkGuiUpdate();

            if (!STATE.isModelChanged) {
                await renderResult();
            }

            rafId = requestAnimationFrame(renderPrediction);
        };

        aiapp().then(e => console.log("loaded data: ", e)).catch(ex => console.log("got exception: ", ex));

        return () => {
            console.log('Running clean-up of effect on unmount')
        }
    });

    const SlideTransition = function (props) {
        return <Slide {...props} direction="up" />;
    }
    let user = appState.userInfo;
    let userInfo = {}
    if (user && user.userInfo) {
        userInfo = JSON.parse(user.userInfo);
    }
    return (

        <div className="App">
            <div id="stats"></div>
            <div id="main">
                <div className="container">
                    <img src={banner} style={{ minWidth: '60vh', marginBottom: '50px' }} />
                    <div className="canvas-wrapper">
                        <canvas id="output"></canvas>
                        <video id="video" ref={video} playsInline className='video'></video>
                    </div>
                    <Card sx={{ minWidth: 640, maxWidth: 640, backgroundColor: user && user.userName ? 'lightgreen' : 'lightgray' }}>
                        {/* <CardMedia
                            sx={{ height: 140 }}
                            image="/static/images/cards/contemplative-reptile.jpg"
                            title="green iguana"
                        /> */}
                        <CardContent>
                            <Typography gutterBottom variant="h5" component="div">
                                User: {user.userName}
                            </Typography>
                            <Typography gutterBottom variant="h5" component="div">
                                Conf Score: {user.score?.toFixed(2)}
                            </Typography>
                            <Grid container spacing={1}>
                                <Grid xs={4}>
                                    <Typography style={{textAlign: 'left'}}>First Name</Typography>
                                </Grid>
                                <Grid xs={8}>
                                    <Typography>{userInfo?.first}</Typography>
                                </Grid>
                                <Grid xs={4}>
                                    <Typography style={{textAlign: 'left'}}>Last Name:</Typography>
                                </Grid>
                                <Grid xs={8}>
                                    <Typography>{userInfo?.last}</Typography>
                                </Grid>
                                <Grid xs={4}>
                                    <Typography style={{textAlign: 'left'}}>Email:</Typography>
                                </Grid>
                                <Grid xs={8}>
                                    <Typography>{userInfo?.email}</Typography>
                                </Grid>
                                <Grid xs={4}>
                                    <Typography style={{textAlign: 'left'}}>Phone Number:</Typography>
                                </Grid>
                                <Grid xs={8}>
                                    <Typography>{userInfo?.phone}</Typography>
                                </Grid>
                               
                            </Grid>
                        </CardContent>

                    </Card>
                    <Snackbar
                        autoHideDuration={6000}
                        TransitionComponent={SlideTransition}
                        anchorOrigin={{ vertical: 'bottom', horizontal: 'right' }}>
                        <MuiAlert severity="info" sx={{ width: '100%' }}>
                            {appState.message}
                        </MuiAlert>
                    </Snackbar>
                </div>
            </div>
        </div>
    );
}

export default App;
