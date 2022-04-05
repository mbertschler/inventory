# Inventory management tool

[![GoDoc](https://godoc.org/github.com/mbertschler/inventory?status.svg)](https://godoc.org/github.com/mbertschler/inventory)
![status: alpha](https://img.shields.io/badge/status-alpha-red.svg)
[![GoDoc](https://goreportcard.com/badge/github.com/mbertschler/inventory)](https://goreportcard.com/report/github.com/mbertschler/inventory)

This tool was designed to simplify the production of electronic hardware. It is used for cataloging newly bought parts and for checking out parts that are used. This tool is already successfully in use, but should still be considered alpha.

## Features

- Mobile optimized web app
- Scans QR part codes with your phone for a quick workflow
  - create new items
  - look up part infos
  - check out parts
- Filter through your parts catalog
- Stores data in a single database file with [bbolt](https://github.com/etcd-io/bbolt)

## Screenshots

<img alt="inventory overview" src="https://mbertschler.com/github/inventory/overview.png" width="480"/><img alt="part detail" src="https://mbertschler.com/github/inventory/detail.png" width="380"/>

## Install

```bash
go install github.com/mbertschler/inventory@latest
```

## Run

Provide the path to your database file. It will be created if it doesn't exist.

```bash
inventory ./inventory.db
```

Once inventory is running, visit http://localhost:5080.

## License

Apache 2.0 with Commons Clause - see [LICENSE](LICENSE)

---

Created by [mbertschler](https://github.com/mbertschler) and [typerat](https://github.com/typerat) in 2018.
