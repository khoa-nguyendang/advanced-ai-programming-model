import json
import logging
import time
from json import JSONEncoder

import cv2
import numpy as np

from .face_det_torch.predict_model import FaceDet
from .face_reg_torch.predict_model import FaceReg
from .mask_cls_torch.predict_model import MaskClassify

logging.basicConfig(format='%(asctime)s,%(msecs)03d %(levelname)-8s [%(filename)s:%(lineno)d] %(message)s',
    datefmt='%Y-%m-%d:%H:%M:%S',
    level=logging.DEBUG)

logger = logging.getLogger(__name__)

_DETECTOR = FaceDet()
_RECOGNIZER = FaceReg()
_MASKCLASSIFIER = MaskClassify()

class NumpyArrayEncoder(JSONEncoder):
    def default(self, obj):
        if isinstance(obj, np.ndarray):
            return obj.tolist()
        return JSONEncoder.default(self, obj)


def extract_feature(image) -> (list or None):
    fvector = []
    try:
        start = time.time()
        logger.info("detect_face")
        all_faces = _DETECTOR.predict(image)
        if len(all_faces) < 1:
            logger.info("unable to find any faces")
            return None

        logger.info("extract_feature: {}".format(all_faces[0][0].shape))
        largest_face_index = 0
        masked = _MASKCLASSIFIER.predict([all_faces[0][0]])
        logger.info("masked: {}".format(masked))
        if masked is not None:
            logger.info("masked[0]: {}".format(masked[0]))
            if masked[0][0] == 1:
                logger.info("Masked detected. Ignore request")
                return None
        logger.info("largest_face_index {}".format(largest_face_index))
        fvector = _RECOGNIZER.predict(all_faces[largest_face_index])
    except Exception as e:
        logger.error(e)
        return None

    logger.info("Vector value size: {}, column size".format(len(fvector)))
    logger.info("Vector value: {}".format(fvector[0]))
    logger.info("Running time {} ".format(time.time() - start))
    return fvector[0]


def _calculate_square_rec(left, top, right, bottom):
    return (right - left) * (bottom - top)


def _test(image_file, dump = True):
    with open(image_file, "rb") as f:
        bytes = f.read(-1)
    npimg = np.frombuffer(bytes, dtype=np.int8)
    img = cv2.imdecode(npimg, cv2.IMREAD_COLOR)

    f = extract_feature(img)
    if f is None:
        return
    if len(f) != 512:
        logger.error("Feature extraction ERROR")
    else:
        logger.info("Feature extraction works!")
        if dump:
            numpyData = {"data": f}
            encodedNumpyData = json.dumps(numpyData, cls=NumpyArrayEncoder)
            logger.info("=========================")
            logger.info(encodedNumpyData)
            logger.info("=========================")


def test_engine():
    # _test("/home/anhcoder/repos/github.com/khoa-nguyendang/advanced-ai-programming-infrastructure/externals/engine-service/test_data/test_face1.jpg", False)
    # _test("/home/anhcoder/repos/github.com/khoa-nguyendang/advanced-ai-programming-infrastructure/externals/engine-service/test_data/test_face2.jpg", False)
    # _test("/home/anhcoder/repos/github.com/khoa-nguyendang/advanced-ai-programming-infrastructure/externals/engine-service/test_data/test_face3.jpg", False)
    _test("/aapi-engine-service/test_data/test_face1.jpg", False)
    _test("/aapi-engine-service/test_data/test_face2.jpg", False)
    _test("/aapi-engine-service/test_data/test_face3.jpg", False)
