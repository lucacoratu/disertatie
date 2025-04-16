import pickle
import sys

import numpy as np
import pandas as pd

# Model import
from sklearn.svm import SVC

#Load the trained model from disk
with open('detection/ai/models/svc.pkl', 'rb') as f:
    loaded_model = pickle.load(f)

#Get the features from the runner
if len(sys.argv) != 2:
    print("Usage: python3 svc.py features.\nEx: python3 svc.py 24,1,0,0.000000,0,0,0,0,0,0,0,0,0,0,0,0,0,0,-1.000000,-1.000000,-1.000000,-1.000000,-1.000000,-1.000000,-1.000000,-1.000000")
    exit(-1)

#Convert the features from string to list of floats
features = [float(feature) for feature in sys.argv[1].split(',')]

#Benign entry
names = ['UrlLength', 'NumberParams', 'NumberSpecialChars', 'RatioSpecialChars', 'NumberRoundBrackets', 'NumberSquareBrackets', 'NumberCurlyBrackets', 'NumberApostrophes', 'NumberQuotationMarks', 'NumberDots', 'NumberSlash', 'NumberBackslash', 'NumberComma', 'NumberColon', 'NumberSemicolon', 'NumberMinus', 'NumberPlus','NumberLessGreater', 'DistanceDots', 'DistanceSlash', 'DistanceBackslash', 'DistanceComma', 'DistanceColon', 'DistanceSemicolon', 'DistanceMinus', 'DistancePlus']

assert len(features) == len(names)

#Reshape the features to the accepted data for the svc model
a = pd.DataFrame(np.array(features).reshape(1, -1))
#Add the names of the columns
a.columns = names

# Make the prediction
prediction = loaded_model.predict(a)

#Return the prediction to the runner
print(prediction[0])