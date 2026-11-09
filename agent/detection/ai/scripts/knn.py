import pickle
import sys

import numpy as np
import pandas as pd

# Model import
from sklearn.neighbors import KNeighborsClassifier

#Names of the features (same order as RequestFeatures.ToString() in the agent)
names = ['UrlLength', 'NumberParams', 'NumberSpecialChars', 'RatioSpecialChars', 'NumberRoundBrackets', 'NumberSquareBrackets', 'NumberCurlyBrackets', 'NumberApostrophes', 'NumberQuotationMarks', 'NumberDots', 'NumberSlash', 'NumberBackslash', 'NumberComma', 'NumberColon', 'NumberSemicolon', 'NumberMinus', 'NumberPlus','NumberLessGreater', 'DistanceDots', 'DistanceSlash', 'DistanceBackslash', 'DistanceComma', 'DistanceColon', 'DistanceSemicolon', 'DistanceMinus', 'DistancePlus']


def classify(model, line):
    #Convert the features from string to list of floats
    features = [float(feature) for feature in line.split(',')]
    if len(features) != len(names):
        raise ValueError("expected %d features, got %d" % (len(names), len(features)))

    #Reshape the features to the accepted data for the model and add the column names
    a = pd.DataFrame(np.array(features).reshape(1, -1), columns=names)

    #Make the prediction
    return model.predict(a)[0]


#Load the trained model from disk (only once, at startup)
with open('detection/ai/models/knn.pkl', 'rb') as f:
    loaded_model = pickle.load(f)

#One-shot mode (for manual testing): python3 knn.py <features>
if len(sys.argv) == 2:
    print(classify(loaded_model, sys.argv[1]))
    sys.exit(0)

#Pipe mode (used by the agent): signal that the model is loaded, then
#read one line of comma separated features from stdin and write one line
#with the prediction to stdout, until stdin is closed
print("READY", flush=True)
for line in sys.stdin:
    line = line.strip()
    if not line:
        continue
    try:
        print(classify(loaded_model, line), flush=True)
    except Exception as e:
        #Keep the request/response stream in sync even on bad input
        print("ERROR " + str(e).replace("\n", " "), flush=True)
