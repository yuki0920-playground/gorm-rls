
# RLSのExample

## セットアップ

ログイン
```
PGPASSWORD=password psql -h localhost -U tenant_user -d postgres -p 5432
```

app.tenant_idをセット
```
postgres=> SET app.tenant_id = 'tenant1';
SET
```

## RLS有効化によるクエリ

### SELECT

tenant_idを指定せずともtenant_idで絞り込まれている
```
postgres=> select * from projects;
    id    |    name     | tenant_id
----------+-------------+-----------
 project1 | Project One | tenant1
 project2 | Project Two | tenant1
(2 rows)
```

自分のtenant_idを指定してもクエリできる

```
postgres=# select * from projects where tenant_id = 'tenant1';
    id    |    name     | tenant_id
----------+-------------+-----------
 project1 | Project One | tenant1
 project2 | Project Two | tenant1
(2 rows)
```

他のtenant_idを指定するとクエリ結果が0件になる

```
postgres=> select * from projects where tenant_id = 'tenant2';
 id | name | tenant_id
----+------+-----------
(0 rows)
```

### INSERT

他のtenant_idを指定するとエラーになる
```
postgres=> INSERT INTO projects (id, name, tenant_id) VALUES ('project5', 'Project One', 'tenant2');
ERROR:  new row violates row-level security policy for table "projects"
```

tenant_idを指定しなくてもエラーになる
```
postgres=> INSERT INTO projects (id, name) VALUES ('project5', 'Project One');
ERROR:  new row violates row-level security policy for table "projects"
```

自分のtenant_idを指定するとINSERTできる
```
postgres=> INSERT INTO projects (id, name, tenant_id) VALUES ('project5', 'Project One', 'tenant1');
INSERT 0 1
```

### UPDATE

他のtenant_idを指定しても更新できない(エラーにはならない)
```
postgres=> update projects set name ='Project One Updated' where tenant_id = 'tenant2' and name = 'Project One';
UPDATE 0
postgres
```

他のtenant_idのレコードを指定しても更新できない(エラーにはならない)
```
postgres=> update projects set name ='Project Three Updated' where name = 'Project Three';
UPDATE 0
```

自分のtenant_idのものなら指定しなくても更新できる
```
postgres=> update projects set name ='Project One Updated' where name = 'Project One';
UPDATE 1
```

自分のtenant_idを指定しても実行できる
```
postgres=> update projects set name ='Project One Updated' where name = 'Project One' and tenant_id = 'tenant1';
UPDATE 1
```


### DELETE

他のtenant_idを指定しても削除できない(エラーにはならない)
```
postgres=> delete from projects where tenant_id = 'tenant2';
DELETE 0
```

自分のtenant_idなら指定しなくても削除できる
```
delete from projects where id = 'project5';
DELETE 1
```

自分のtenant_idを指定しても削除できる
```
postgres=# delete from projects where id = 'project5' and tenant_id = 'tenant1';
DELETE 1
```
