# Author: Nguyen Y Hop
import glob
import json
import os

import cv2
import numpy as np
import torch

from .src.models.architectures.base_model import RetinaFace
from .src.prior_box import PriorBox
from .src.utils import decode, decode_landm, py_cpu_nms


class FaceDet:

    def __init__(self, **kwargs):
        
        self.dynamic_path = os.path.dirname(os.path.realpath(__file__))
        self.config = self.load_config(os.path.join(self.dynamic_path, 'configs/face_config.json'))
       
        self.prior_box_config = self.config['priorbox_cfg']

        self.height, self.width = self.config['image_size'][:2]
        self.priorbox = PriorBox(self.prior_box_config)
        
        self.resize = self.config['resize']
        self.confidence_threshold = self.config['confidence_threshold']
        self.top_k = self.config['top_k']
        self.variance = self.config['variance']
        self.keep_top_k = self.config['keep_top_k']
        self.nms_threshold = self.config['nms_threshold']
        self.vis_thres = self.config['vis_thres']

   
        model_path = os.path.join(self.dynamic_path, self.config['model_path'])
        

        if os.path.exists(model_path):
            self.build_model(model_path) 
        else:
            self.model = None
            print("Cannot find the weight's path")
        print("DONE")
        

    def remove_prefix(self, state_dict, prefix):
        ''' Old style model is stored with all names of parameters sharing common prefix 'module.' '''
        print('remove prefix \'{}\''.format(prefix))
        f = lambda x: x.split(prefix, 1)[-1] if x.startswith(prefix) else x
        return {f(key): value for key, value in state_dict.items()}

    def check_keys(self, model, pretrained_state_dict):
        ckpt_keys = set(pretrained_state_dict.keys())
        model_keys = set(model.state_dict().keys())
        used_pretrained_keys = model_keys & ckpt_keys
        unused_pretrained_keys = ckpt_keys - model_keys
        missing_keys = model_keys - ckpt_keys
        print('Missing keys:{}'.format(len(missing_keys)))
        print('Unused checkpoint keys:{}'.format(len(unused_pretrained_keys)))
        print('Used keys:{}'.format(len(used_pretrained_keys)))
        assert len(used_pretrained_keys) > 0, 'load NONE from pretrained checkpoint'
        return True


    def build_model(self, model_path):
        self.model = RetinaFace(self.config['architecture'])
        print('Loading pretrained model from {}'.format(model_path))
        if self.config['load_to_cpu']:
            pretrained_dict = torch.load(model_path, map_location=lambda storage, loc: storage)
        else:
            device = torch.cuda.current_device()
            pretrained_dict = torch.load(model_path, map_location=lambda storage, loc: storage.cuda(device))
        if "state_dict" in pretrained_dict.keys():
            pretrained_dict = self.remove_prefix(pretrained_dict['state_dict'], 'module.')
        else:
            pretrained_dict = self.remove_prefix(pretrained_dict, 'module.')
        self.check_keys(self.model, pretrained_dict)
        self.model.load_state_dict(pretrained_dict, strict=False)
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


    def pre_process(self, img_raw):
        img = np.float32(img_raw)
        ori_h, ori_w = img.shape[:2]
        scale = np.array([ori_w, ori_h, ori_w, ori_h])
        img -= (104, 117, 123)
        img = img.transpose(2, 0, 1)
        img = torch.from_numpy(img).unsqueeze(0)
        return img, scale


    def visual(self, img_raw, dets):
        for b in dets:
            if b[4] < self.vis_thres:
                continue
            text = "{:.4f}".format(b[4])
            b = list(map(int, b))
            
            cv2.rectangle(img_raw, (b[0], b[1]), (b[2], b[3]), (0, 0, 255), 2)
            cx = b[0]
            cy = b[1] + 12
            cv2.putText(img_raw, text, (cx, cy),
                        cv2.FONT_HERSHEY_DUPLEX, 0.5, (255, 255, 255))

            # landms
            cv2.circle(img_raw, (b[5], b[6]), 1, (0, 0, 255), 4)
            cv2.circle(img_raw, (b[7], b[8]), 1, (0, 255, 255), 4)
            cv2.circle(img_raw, (b[9], b[10]), 1, (255, 0, 255), 4)
            cv2.circle(img_raw, (b[11], b[12]), 1, (0, 255, 0), 4)
            cv2.circle(img_raw, (b[13], b[14]), 1, (255, 0, 0), 4)
        return img_raw


    def get_face_and_confi(self, img_raw, dets):
        outputs = []
        for b in dets:
            if b[4] < self.vis_thres:
                continue
            xmin, ymin, xmax, ymax, confi = list(map(int, b[:5]))
            face = img_raw[ymin:ymax, xmin:xmax]
            if min(face.shape[:2]) < 5:
                continue
            outputs.append([face, round(confi, 4)])
        return outputs


    @torch.no_grad()
    def predict(self, img, visual_mode=False):
        
        img_raw = img.copy()
        ori_h, ori_w = img.shape[:2]
        img, scale = self.pre_process(img)
        loc, conf, landms = self.model(img, training=False)
        prior_data = self.priorbox.forward(image_size=[ori_h, ori_w])
    

        boxes = decode(loc.data.squeeze(0), prior_data, self.variance)
        boxes = boxes * scale / self.resize
        boxes = boxes.numpy()
        scores = conf.squeeze(0).data.numpy()[:, 1]
    
        landms = decode_landm(landms.data.squeeze(0), prior_data, self.variance).numpy()
        scale1 = np.array([ori_w, ori_h, ori_w, ori_h,
                                ori_w, ori_h, ori_w, ori_h,
                                ori_w, ori_h])
        landms = landms * scale1 / self.resize
        landms = landms

        inds = np.where(scores > self.confidence_threshold)[0]
        boxes = boxes[inds]
        landms = landms[inds]
        scores = scores[inds]

        # keep top-K before NMS
        order = scores.argsort()[::-1][:self.top_k]
        boxes = boxes[order]
        landms = landms[order]
        scores = scores[order]

        dets = np.hstack((boxes, scores[:, np.newaxis])).astype(np.float32, copy=False)
        keep = py_cpu_nms(dets, self.nms_threshold)
        # keep = nms(dets, args.nms_threshold,force_cpu=args.cpu)
        dets = dets[keep, :]
        landms = landms[keep]

        # keep top-K faster NMS
        dets = dets[:self.keep_top_k, :]
        landms = landms[:self.keep_top_k, :]

        dets = np.concatenate((dets, landms), axis=1)

        if visual_mode:
            img_raw = self.visual(img_raw, dets)
            path = "/home/anhcoder/repos/github.com/khoa-nguyendang/advanced-ai-programming-infrastructure/externals/research/assets/outputs"
            cv2.imwrite(os.path.join(path , 'output.jpg'), img_raw)
            return img_raw
        outputs = self.get_face_and_confi(img_raw, dets)
        return outputs


if __name__ == '__main__':
    import time
    face_det = FaceDet()
    for img_path in glob.glob("DATATEST/image/*.*"):
        img_raw = cv2.imread(img_path, cv2.IMREAD_COLOR)
        img = np.float32(img_raw)
        start = time.time()
        img_pred = face_det.predict(img, visual_mode=False)
        print("Time:\t", time.time() - start)
        for (img, score) in img_pred:
            cv2.imwrite(f"outputs/module/{time.time()}.jpg", img)