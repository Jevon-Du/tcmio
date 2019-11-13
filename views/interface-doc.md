
# TCMIO v1.0 Interface document

## Request Url  

Url                                     |   Method  |   Parametes
:---                                    |   :---    |   :---
/targets                                |   get     |   draw=1&start=0&length=5&order[0][column]=0&order[0][dir]=asc
/targets/:id([0-9]+)                    |   get     |   id
targets/:id([0-9]+)/json                |   get     |   id
/ligands                                |   get     |   draw=1&start=0&length=5&order[0][column]=0&order[0][dir]=asc
/ligands/:id([0-9]+)                    |   get     |   id
ligands/:id([0-9]+)/json                |   get     |   id
/ingredients                            |   get     |   draw=1&start=0&length=5&order[0][column]=0&order[0][dir]=asc
/ingredients/:id([0-9]+)                |   get     |   id
ingredients/:id([0-9]+)/json            |   get     |   id
/tcms                                   |   get     |   draw=1&start=0&length=5&order[0][column]=0&order[0][dir]=asc
/tcms/:id([0-9]+)                       |   get     |   id
tcms/:id([0-9]+)/json                   |   get     |   id
/prescriptions                          |   get     |   draw=1&start=0&length=5&order[0][column]=0&order[0][dir]=asc
/prescriptions/:id([0-9]+)              |   get     |   id
prescriptions/:id([0-9]+)/json          |   get     |   id
network/tcms                            |   get     |   kw=丁香,九里香&type=chinese_name
network/prescriptions                   |   get     |   kw=Yi Qing Jiao Nang&type=pinyin_name
structure/ligand                        |   get     |   query=Molecule%20from%20ipmDraw&method=sim&threshold=0.9&type=ligand
structure/ingredient                    |   get     |   query=Molecule%20from%20ipmDraw&method=sim&threshold=0.9&type=ingredient



## Response

### /targets, /ligands, /ingredients, /tcms, /prescriptions

```json
{
    "draw": 1,
    "recordsTotal" : 497,
    "recordsFiltered" : 100,
    "data": [{…}, …]
}
```

### targets/:id([0-9]+)/json & …

```json
{
    "Data": {"Id": 8, "Name": "Carbonic anhydrase 9 ", "GeneName": "CA9",…},
    "Msg": "",
    "State": "success"
}
```

### network/tcms & network/prescriptions

```json
{
    "Data": {"edges":[{…},…], "nodes":[{…},…],
    "Msg": "",
    "State": "success"
}
```

### structure/ligand & structure/ingredient

```json
{
    "Data": ,
    "Msg": "",
    "State": "success"
}
```

#### 功能说明

#####　前置条件

进入MOA页面

#### 功能描述

1. 显示图网络
2. 点击节点展示该节点的简要信息

####　业务规则

1. 节点分为5种, prescription, tcm, ingredient, liagnd
2. 各节点信息同browse页面的信息
3. pathway?


TODO:

[] 每个节点制作不同的logo
[] Download Figure
[] Download Network ??
[] pathway的数据请求,如何触发

