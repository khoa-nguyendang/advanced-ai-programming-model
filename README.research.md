## Install lib
```
pip install torchvison, cv2, numpy
```

## Inference face detect module

    ▶️ Input: <np.array> a color image
<img src="research/assets/inputs.jpg" style=" width:640px ; height:360px "  >

    ▶️ Ouptut: <list> face images and score list (list)

    [[face_img1, 0.98], [face_img2, 0.97]...]

<img src="research/assets/output1.jpg" style=" width:96px ; height:96px "  >
<img src="research/assets/output2.jpg" style=" width:96px ; height:96px "  >
<img src="research/assets/output3.jpg" style=" width:96px ; height:96px "  >
<img src="research/assets/output4.jpg" style=" width:96px ; height:96px "  >



### Script

```
from research.face_det_torch.predict_model import FaceDet

face_det = FaceDet()

img = cv2.imread(img_path)
img_pred = face_det.predict(img)

```


## Inference face reg module

    ▶️ Input: <list> face images list from face detection module

        [face_img1, face_img2, ...]

    ▶️ Ouptut: <list> visual 512 dim vector of face respectively

        [vector_face_1, vector_face_2 ...]



### Script

```
from research.face_reg_torch.predict_model FaceReg

face_reg = FaceReg()

vector = face_reg.predict(img_list)

```