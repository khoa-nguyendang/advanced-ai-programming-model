import torch
from torch import nn


class FacemaskRecognizeModel(nn.Sequential):
    def __init__(self):
        super(FacemaskRecognizeModel, self).__init__(
            # first layer
            nn.Conv2d(in_channels = 3, out_channels = 6, kernel_size = 5),
            nn.ReLU(),
            nn.MaxPool2d(2, 2),

            # second layer
            nn.Conv2d(in_channels = 6, out_channels = 12, kernel_size = 5),
            nn.ReLU(),
            nn.MaxPool2d(2, 2),

            # flatten
            nn.Flatten(),

            # linear layers
            nn.Linear(in_features = 300, out_features = 150),
            nn.Linear(in_features = 150, out_features = 1),

            # active function
            nn.Sigmoid()
        )
