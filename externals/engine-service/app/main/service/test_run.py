
import time

import cv2
from face_det_torch.predict_model import FaceDet
from face_reg_torch.predict_model import FaceReg

start = time.time()
face_det = FaceDet()
img_path = "/home/anhcoder/repos/github.com/khoa-nguyendang/advanced-ai-programming-infrastructure/externals/engine-service/test_data/inputs/picture(1).jpg"
img = cv2.imread(img_path)
img_pred = face_det.predict(img)

face_reg = FaceReg()

vector = face_reg.predict(img_pred)

end = time.time()
print(end - start)
print(len(vector))
print('done \n')