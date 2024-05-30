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

Next, perform aggregation to result csv file with interesting TCP labels and packet states.
pcapexporter catches TCP retransmissions (fast, spurious and simple retransmissions), and RST(reset) packet. See `pkg/analyzer/crawler/crawler.go:matchData`

```shell
./bin/pcapexporter analyze -p <PCAP-FILES-DIR> -o <RES>
```

pcapexporter will scan all files with mask `*.pcap*` (to be compatible with wireshark file rotation mechanism) and migrate them to csv with `tshark`

Analyze performed in parallel with 20 goroutines by default. You can specify number of goroutines with `--workers`

Result csv can be imported into db or analyzed manually

### Further example

Next, up grafana + mysql

```shell
docker compose up -d
```

Prepare MySQL

```shell
# Create database and schema
docker exec -i pcap-mysql mysql -uroot pcap < mysql/schema.sql
# Grant access to grafana user
docker exec -i pcap-mysql mysql -uroot pcap < mysql/access.sql
```

Import csv:
Run commands from `mysql/mysqlimport.bash` to import `record.csv`

Fill `record_grouped_rt` table

```shell
docker exec -i pcap-mysql mysql -uroot pcap < mysql/refill.sql
```

Go to `localhost:3000`. Import Dashboard from `grafana/dashboad.json`


