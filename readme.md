#go-mysql-rabbitmq
Golang реализация протокола MySQL репликации.
На каждый SQL запрос изменяющий данные в MySQL отправляется сообщение в RabbitMQ.

Сообщения попадают в Exchange с именем заданным в парметре `-rabbitmq_exchange` и routing key состоящим из  `[schema].[table name].[event type]`



##Виды сообщений

* insert
* update
* delete

####Insert

```
{
  "Table": {
    "Schema": "social",
    "Name": "user",
    "Columns": [
      {
        "Name": "id",
        "Type": 1,
        "RawType": "int(11)",
        "IsAuto": true,
        "EnumValues": null,
        "SetValues": null
      },
      {
        "Name": "name",
        "Type": 5,
        "RawType": "varchar(255)",
        "IsAuto": false,
        "EnumValues": null,
        "SetValues": null
      }
    ],
    "Indexes": [
      {
        "Name": "PRIMARY",
        "Columns": [
          "id"
        ],
        "Cardinality": [
          1
        ]
      }
    ],
    "PKColumns": [
      0
    ]
  },
  "Action": "insert",
  "Rows": [
    [
      1,
      "John"
    ]
  ]
}

```

####Update
```
{
  "Table": {
    "Schema": "social",
    "Name": "user",
    "Columns": [
      {
        "Name": "id",
        "Type": 1,
        "RawType": "int(11)",
        "IsAuto": true,
        "EnumValues": null,
        "SetValues": null
      },
      {
        "Name": "name",
        "Type": 5,
        "RawType": "varchar(255)",
        "IsAuto": false,
        "EnumValues": null,
        "SetValues": null
      }
    ],
    "Indexes": [
      {
        "Name": "PRIMARY",
        "Columns": [
          "id"
        ],
        "Cardinality": [
          1
        ]
      }
    ],
    "PKColumns": [
      0
    ]
  },
  "Action": "update",
  "Rows": [
    [
      1,
      "John"
    ],
    [
      1,
      "Doe"
    ]
  ]
}

```

####Delete

```
{
  "Table": {
    "Schema": "social",
    "Name": "user",
    "Columns": [
      {
        "Name": "id",
        "Type": 1,
        "RawType": "int(11)",
        "IsAuto": true,
        "EnumValues": null,
        "SetValues": null
      },
      {
        "Name": "name",
        "Type": 5,
        "RawType": "varchar(255)",
        "IsAuto": false,
        "EnumValues": null,
        "SetValues": null
      }
    ],
    "Indexes": [
      {
        "Name": "PRIMARY",
        "Columns": [
          "id"
        ],
        "Cardinality": [
          1
        ]
      }
    ],
    "PKColumns": [
      0
    ]
  },
  "Action": "delete",
  "Rows": [
    [
      1,
      "Doe"
    ]
  ]
}

```

## Параметры для запуска go-mysql-rabbitmq

```
  -data-dir string
        Path to store data, like master.info (default "./var")
  -flavor string
        Flavor: mysql or mariadb (default "mysql")
  -host string
        MySQL host (default "127.0.0.1")
  -mysqldump string
        mysqldump execution path
  -password string
        MySQL password
  -port int
        MySQL port (default 3306)
  -rabbitmq_exchange string
        RabbitMQ exchange name (default "mysql")
  -rabbitmq_host string
        RabbitMQ host (default "127.0.0.1")
  -rabbitmq_password string
        RabbitMQ password (default "guest")
  -rabbitmq_port int
        RabbitMQ port (default 5672)
  -rabbitmq_user string
        RabbitMQ user (default "guest")
  -server-id int
        Unique Server ID (default 101)
  -user string
        MySQL user, must have replication privilege (default "root")
```