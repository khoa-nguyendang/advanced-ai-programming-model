import json
import os

import cv2
import numpy as np
import torch
from PIL import Image
from torchvision import datasets, transforms

from .src.classifcation_model import FacemaskRecognizeModel


class MaskClassify:

    def __init__(self, **kwargs):
        self.dynamic_path = os.path.dirname(os.path.realpath(__file__))
        self.config = self.load_config(os.path.join(self.dynamic_path, 'configs/maskcls_config.json'))
        model_path = os.path.join(self.dynamic_path, self.config['model_path'])

        self.transform = transforms.Compose([
                transforms.Resize((32,32)),
                transforms.ToTensor()
            ])
    
        if os.path.exists(model_path):
            self.build_model(model_path) 
        else:
            self.model = None
            print("Cannot find the weight's path")
        print("DONE")

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


    
    def build_model(self, model_path):
        self.model = FacemaskRecognizeModel()
        print('Loading pretrained model from {}'.format(model_path))
        use_gpu = torch.cuda.is_available()
        print("USE GPU:\t", use_gpu)
        self.device = torch.device("cpu")
        if not use_gpu:
            print("LOAD INTO CPU")
            pretrained_dict = torch.load(model_path, map_location=lambda storage, loc: storage)
        else:
            print("LOAD INTO GPU")
            # device = torch.cuda.current_device()
            self.device = torch.device("cuda")
            pretrained_dict = torch.load(model_path, map_location=lambda storage, loc: storage.cuda(self.device))
       
        self.model.load_state_dict(pretrained_dict, strict=False)
        self.model.to(self.device)
        self.model.eval()

    
    def preproces(self, img_list):
        batch = []
        for img in img_list:
            if not isinstance(img, np.ndarray):
                continue
            img = cv2.cvtColor(img, cv2.COLOR_BGR2RGB)
            img = Image.fromarray(img)
            img = self.transform(img)
            batch.append(img[None])
        if len(batch) == 0:
            return None
        return torch.concat(batch, axis=0)


    def predict(self, imgs):
        batch = self.preproces(imgs)
        if batch is None:
            return []
        prob = self.model(batch)
        predicts = torch.round(prob).detach().numpy()
        cls_preds = np.round(predicts)
        outputs = []
        for (cls_, score) in zip(cls_preds, predicts):
            score = score if cls_ else 1- score
            outputs.append([cls_[0],score[0]])
        return outputs
      
    


if __name__ == '__main__':
   
    
    class_ = MaskClassify()
    img1 = cv2.imread('debugs/mask/0001.png')
    img2 = cv2.imread('debugs/mask/0002.png')
    outputs = class_.predict([img1, img2])
    print(outputs)