# pcapexporter

Exports pcap analyze to csv file.

Uses:
* tshark - to migrate pcap file to csv
* `github.com/google/gopacket` - to parse pcap packet from Go
* CGO - to use `github.com/google/gopacket`

## Usage

Build pcapexporter

```shell
export CGO_ENABLED=1
go build -v -o ./bin/pcapexporter ./cmd/pcapexporter 
```

Perform pcap file migration to csv:
```shell
./bin/pcapexporter csv --path <PCAP-FILES-DIR>
```
pcapexporter will scan all files with mask `*.pcap*` (to be compatible with wireshark file rotation mechanism) and migrate them to csv with `tshark`

Next, perform aggregation to result csv file with interesting TCP labels and packet states.
pcapexporter catches TCP retransmissions (fast, spurious and simple retransmissions), and RST(reset) packet. See `pkg/analyzer/crawler/crawler.go:matchData`
```shell
./bin/pcapexporter analyze --path <PCAP-FILES-DIR>
```
Analyze performed in parallel with 20 goroutines by default. You can specify number of goroutines with `--workers`

Result file `data.csv` will be located in current directory.

`data.csv` can be imported into db or analyzed manually

### Further example

Next, up grafana + mysql

```shell
docker compose up -d
```

Prepare MySQL

```shell
# Create database and schema
docker compose mysql exec -i mysql -uroot < mysql/schema.sql
# Grant access to grafana user
docker compose mysql exec -i mysql -uroot < mysql/access.sql
```

Import csv:
Run commands from `mysql/mysqlimport.bash` to import `record.csv`

Fill `record_grouped_rt` table

```shell
docker compose mysql exec -i mysql -uroot < mysql/refill.sql
```


Go to `localhost:3000`. Import Dashboard from `grafana/dashboad.json`


