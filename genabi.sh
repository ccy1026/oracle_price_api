cd utils/contract/abi;
for FILE in *;
  do echo ${FILE%%.*} success && abigen --abi=./$FILE --pkg=contract --type=${FILE%%.*} --out=../${FILE%%.*}.go;
done
