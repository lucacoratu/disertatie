import pandas as pd

# Modelling
from sklearn.preprocessing import StandardScaler
from sklearn.svm import SVC
from sklearn.metrics import accuracy_score, confusion_matrix, precision_score, recall_score, classification_report
from sklearn.model_selection import train_test_split

names = ['UrlLength', 'NumberParams', 'NumberSpecialChars', 'NumberRoundBrackets', 'NumberSquareBrackets', 'NumberCurlyBrackets', 'NumberApostrophes', 'NumberQuotationMarks', 'NumberDots', 'NumberSlash', 'NumberBackslash', 'NumberComma', 'NumberColon', 'NumberSemicolon', 'NumberMinus', 'NumberPlus', 'DistanceDots', 'DistanceSlash', 'DistanceBackslash', 'DistanceComma', 'DistanceColon', 'DistanceSemicolon', 'DistanceMinus', 'DistancePlus']

#Read the csv files
df_benign = pd.read_csv("../../agent/datasets/benign.csv", header=None, names=names)
# df_lfi1 = pd.read_csv("../../agent/datasets/lfi_jhaddix.csv", header=None, names=names)
# df_lfi2 = pd.read_csv("../../agent/datasets/lfi_linux.csv", header=None, names=names)
# df_lfi3 = pd.read_csv("../../agent/datasets/lfi_windows.csv", header=None, names=names)

df_lfi = pd.read_csv("../../agent/datasets/lfi.csv", header=None, names=names)

df_sqli = pd.read_csv("../../agent/datasets/sqli.csv", header=None, names=names)

df_test = pd.read_csv("../../agent/datasets/lfi_test.csv", header=None, names=names)
df_sqli_test = pd.read_csv("../../agent/datasets/sqli_test.csv", header=None, names=names)


df_benign['label'] = 'benign'
# df_lfi1['label'] = 'lfi'
# df_lfi2['label'] = 'lfi'
# df_lfi3['label'] = 'lfi'

df_lfi['label'] = 'lfi'
df_sqli['label'] = 'sqli'

df_test['label'] = 'lfi'
df_sqli_test['label'] = 'sqli'

df_test = pd.concat([df_test, df_sqli_test], ignore_index=True)
# print(df_test)

#Concatenate all the datasets
full_df = pd.concat([df_benign, df_lfi, df_sqli] , ignore_index=True)

#Randomize the rows
full_df = full_df.sample(frac=1, random_state=42).reset_index(drop=True)

#Extract the features
X = full_df.drop('label', axis=1)
y = full_df['label']

# scaler = StandardScaler()
# X_scaled = scaler.fit_transform(X)

#Prepare the train data and test data
X_train, X_test, y_train, y_test = train_test_split(
    X, y, test_size=0.2, random_state=42
)

svm = SVC(kernel='rbf', C=1, gamma='scale')  # kernel can be 'linear', 'poly', 'rbf', 'sigmoid'
svm.fit(X_train, y_train)

# X_test_1 = scaler.fit_transform(df_test.drop('label', axis=1))
# print(X_test_1)

#Predict
# X_test = pd.concat([X_test, df_test.drop('label', axis=1)], ignore_index=True)
# y_test = pd.concat([y_test, df_test['label']], ignore_index=True)

y_pred = svm.predict(X_test)

print("Accuracy:", accuracy_score(y_test, y_pred))
print("\nClassification Report:\n", classification_report(y_test, y_pred))


y_pred2 = svm.predict(df_test.drop('label', axis=1))

print("Accuracy:", accuracy_score(df_test['label'], y_pred2))
print("\nClassification Report:\n", classification_report(df_test['label'], y_pred2))
