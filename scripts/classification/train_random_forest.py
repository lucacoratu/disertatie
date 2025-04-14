import pandas as pd
import numpy as np

# Modelling
from sklearn.ensemble import RandomForestClassifier
from sklearn.metrics import accuracy_score, confusion_matrix, precision_score, recall_score, ConfusionMatrixDisplay
from sklearn.model_selection import RandomizedSearchCV, train_test_split

#Read the csv files
df_benign = pd.read_csv("../../agent/datasets/benign.csv")
df_lfi1 = pd.read_csv("../../agent/datasets/lfi_jhaddix.csv")
df_lfi2 = pd.read_csv("../../agent/datasets/lfi_linux_files.csv")
df_lfi3 = pd.read_csv("../../agent/datasets/lfi_windows.csv")


#df_benign.assign(Name='Attack')
df_benign['Attack'] = 'benign'
#df_lfi1.assign(Name='Attack')
df_lfi1['Attack'] = 'lfi'
#df_lfi2.assign(Name='Attack')
df_lfi2['Attack'] = 'lfi'
#df_lfi3.assign(Name='Attack')
df_lfi3['Attack'] = 'lfi'