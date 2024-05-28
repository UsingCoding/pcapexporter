mysqlimport \
  --ignore-lines=1 \
  --fields-terminated-by=, \
  --local -u root \
  pcap \
  record.csv
