# Author: Nguyen Y hop
import json
import os

import cv2
import numpy as np
import torch

from .src.base_model import MobileFaceNet


class FaceReg:

    def __init__(self, device=-1, **kwargs):
        self.dynamic_path = os.path.dirname(os.path.realpath(__file__))
        self.config = self.load_config(os.path.join(self.dynamic_path, 'configs/face_config.json'))
        os.environ["CUDA_VISIBLE_DEVICES"]=str(device)

        self.height, self.width = self.config['image_size'][:2]
        model_params = self.config['model_params']
        weight_path = os.path.join(self.dynamic_path, self.config['model_path'])

        self.model = MobileFaceNet(**model_params)
        self.model.load_state_dict(torch.load(weight_path))
        self.model.eval()


    def load_config(self, config_path):
        '''
        Load config
        args:
            config_path <string>
        '''
        assert os.path.exists(config_path), f"Cannot find {config_path} path"
        with open(config_path, 'r') as rf:
            config = json.load(rf)
            return config

    def pre_process(self, img_list):
        batch_img = []
        for img in img_list:
            if not isinstance(img, np.ndarray):
                continue
            img = cv2.resize(img, (self.width, self.height))
            img = cv2.cvtColor(img, cv2.COLOR_BGR2RGB)
            batch_img.append(img)
        if len(batch_img) == 0:
            return [], False
        batch_img = torch.from_numpy(np.array(batch_img)).float()
        batch_img = np.transpose(batch_img, (0, 3, 1, 2))
        batch_img.div_(255).sub_(0.5).div_(0.5)
        return batch_img, True

    @torch.no_grad()
    def predict(self, img_list):
        batch_img, flag = self.pre_process(img_list)
        if not flag:
            return []
        feat = self.model(batch_img).numpy().tolist()
        return feat


if __name__ == '__main__':
    import glob

    from numpy import dot
    from numpy.linalg import norm

    face_reg = FaceReg()
    batch = []
    for img_path in glob.glob("DATATEST/images/*.*"):
        print(img_path)
        img = cv2.imread(img_path)
        batch.append(img)
    feat = face_reg.predict(batch)
    print(feat)
    cos_sim = dot(feat[1], feat[2])/(norm(feat[1])*norm(feat[2]))
    print(cos_sim)
