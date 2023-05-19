import AccountCircle from '@mui/icons-material/AccountCircle';
import AssignmentIndIcon from '@mui/icons-material/AssignmentInd';
import EmailIcon from '@mui/icons-material/Email';
import PeopleIcon from '@mui/icons-material/People';
import PhoneAndroidIcon from '@mui/icons-material/PhoneAndroid';
import { Button, Paper } from '@mui/material';
import Box from '@mui/material/Box';
import FormControl from '@mui/material/FormControl';
import Input from '@mui/material/Input';
import InputAdornment from '@mui/material/InputAdornment';
import InputLabel from '@mui/material/InputLabel';
import * as React from 'react';

const style = {
    position: 'absolute',
    top: '50%',
    left: '50%',
    transform: 'translate(-50%, -50%)',
    width: 800,
    bgcolor: 'background.paper',
    border: '2px solid #000',
    boxShadow: 24,
    p: 4,
};

export default function AddUser({ onComplete }) {
    const [appState, setAppState] = React.useState({
        userId: '',
        referenceId: '',
        firstName: '',
        lastName: '',
        phone: '',
        email: '',
        userName: '',
        images: [],
    });
    function onSubmit(e) {
        e?.preventDefault();
        // const formData = new FormData();
        // formData.append("UserId", state.UserId);
        // formData.append("ReferenceId", state.ReferenceId);
        let userInfo = {
            first: appState.firstName,
            last: appState.lastName,
            phone: appState.phone,
            email: appState.email,
        }
        // formData.append("UserInfo", JSON.stringify(userInfo));
        // formData.append("UserName", state.UserName);
        // formData.append("Images", state.UserName);
        let formData = {
            userId: appState.userId,
            referenceId: appState.referenceId,
            userName: appState.userName,
            images: appState.images,
            userInfo: JSON.stringify(userInfo),
        }
        fetch('http://localhost:8081/api/v1/user/enroll', {
            method: 'POST',
            body: JSON.stringify(formData),
            headers: {
                'Content-Type': 'application/json'
            },
        })
            .then(res => res.json())
            .then(data => {
                console.log(" enroll result: ", data);
                onComplete(data);
            })
            .catch(ex => console.log("got exception  when enroll, ", ex));
    }

    const handleFileUpload = (e) => {
        if (!e.target.files) {
            return;
        }
        const files = e.target.files;

        for (let i = 0; i < files.length; i++) {
            let file = files[i];
            const reader = new FileReader();
            reader.onload = (evt) => {
                if (!evt?.target?.result) {
                    return;
                }
                
                setAppState(pre => ({
                    ...pre,
                    images: [...pre.images, 
                        {
                            imageId: file.name,
                            data: evt.target.result?.replace("data:image/jpeg;base64,", "")?.replace("data:image/png;base64,", "")
                        }]
                }))
            };
            reader.readAsDataURL(file);
        }
    };
    const onHandleChange = (e) => {
        const { name, value } = e.target;
        setAppState(pre => ({
            ...pre,
            [name]: value
        }))
    }
    return (
        <Box sx={style}>
            <Paper>Add New User</Paper>
            <FormControl style={{ margin: 10 }}>
                <InputLabel htmlFor="userId">
                    User Id
                </InputLabel>
                <Input id="userId" name='userId' value={appState.UserId} onChange={onHandleChange}
                    startAdornment={
                        <InputAdornment position="start">
                            <AssignmentIndIcon />
                        </InputAdornment>
                    }
                />
            </FormControl>
            <FormControl style={{ margin: 10 }}>
                <InputLabel htmlFor="referenceId">
                    Reference Id
                </InputLabel>
                <Input id="referenceId" name='referenceId' value={appState.ReferenceId} onChange={onHandleChange}
                    startAdornment={
                        <InputAdornment position="start">
                            <PeopleIcon />
                        </InputAdornment>
                    }
                />
            </FormControl>
            <FormControl style={{ margin: 10 }}>
                <InputLabel htmlFor="userName">
                    User name
                </InputLabel>
                <Input id="userName" name='userName' value={appState.userName} onChange={onHandleChange}
                    startAdornment={
                        <InputAdornment position="start">
                            <AccountCircle />
                        </InputAdornment>
                    }
                />
            </FormControl>
            <FormControl style={{ margin: 10 }}>
                <InputLabel htmlFor="firstName">
                    First name
                </InputLabel>
                <Input id="firstName" name='firstName' value={appState.firstName} onChange={onHandleChange}
                    startAdornment={
                        <InputAdornment position="start">
                            <AccountCircle />
                        </InputAdornment>
                    }
                />

            </FormControl>
            <FormControl style={{ margin: 10 }}>
                <InputLabel htmlFor="input-with-icon-adornment">
                    Last name
                </InputLabel>
                <Input id="lastName" name='lastName' value={appState.lastName} onChange={onHandleChange}
                    startAdornment={
                        <InputAdornment position="start">
                            <AccountCircle />
                        </InputAdornment>
                    }
                />
            </FormControl>
            <FormControl style={{ margin: 10 }}>
                <InputLabel htmlFor="phone">
                    Phone
                </InputLabel>
                <Input id="phone" name='phone' value={appState.phone} onChange={onHandleChange}
                    startAdornment={
                        <InputAdornment position="start">
                            <PhoneAndroidIcon />
                        </InputAdornment>
                    }
                />

            </FormControl>
            <FormControl style={{ margin: 10 }}>
                <InputLabel htmlFor="email">
                    Email
                </InputLabel>
                <Input id="email" name='email' value={appState.email} onChange={onHandleChange}
                    startAdornment={
                        <InputAdornment position="start">
                            <EmailIcon />
                        </InputAdornment>
                    }
                />

            </FormControl>
            <FormControl style={{ margin: 10 }}>
                <input
                    accept="image/*"
                    style={{ display: 'none' }}
                    id="raised-button-file"
                    multiple
                    type="file"
                    name='images'
                    onChange={handleFileUpload}
                />
                <label htmlFor="raised-button-file">
                    <Button variant="raised" component="span" >
                        Upload Images
                    </Button>
                </label>
            </FormControl>
            <Button onClick={onSubmit}>Submit</Button>
        </Box>
    );
}
