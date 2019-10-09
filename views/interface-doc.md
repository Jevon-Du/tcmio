
# 接口文档

## request
url: 'targets',
type: 'get',
data: {
    offset: 0,
    limit: RECORDS
}

## response
```json
Data: Array(10)
    0: {ChemblId: "CHEMBL5957;", EcNumber: "3.1.3.5", Function: "Hydrolyzes extracellular nucleotides into membrane…sidase activities. {ECO:0000269|PubMed:21933152}.", GeneName: "NT5E", Kegg: "hsa:4907;", …}
    1: {ChemblId: "CHEMBL1772928;", EcNumber: "", Function: "Drug efflux transporter present in a number of ste…. Specifically present in limbal stem cells, wher", GeneName: "ABCB5", Kegg: "hsa:340273;", …}
    2: {ChemblId: "CHEMBL3004;", EcNumber: "7.6.2.2; 7.6.2.3", Function: "Mediates export of organic anions and drugs from t…hione conjugates, leukotriene C4, estradiol-17-be", GeneName: "ABCC1", Kegg: "hsa:4363;", …}
    3: {ChemblId: "CHEMBL5918;", EcNumber: "", Function: "May act as an inducible transporter in the biliary…tatic hepatocytes (By similarity). {ECO:0000250}.", GeneName: "ABCC3", Kegg: "hsa:8714;", …}
    4: {ChemblId: "CHEMBL1795197;", EcNumber: "4.6.1.2", Function: "Receptor for the E.coli heat-stable enterotoxin (E…the endogenous peptides guanylin and uroguanylin.", GeneName: "GUCY2C", Kegg: "hsa:2984;", …}
    5: {ChemblId: "CHEMBL4660;", EcNumber: "3.2.2.6; 2.4.99.20", Function: "Synthesizes the second messengers cyclic ADP-ribos…o moonlights as a receptor in cells of the immune", GeneName: "CD38", Kegg: "hsa:952;", …}
    6: {ChemblId: "CHEMBL3712864;", EcNumber: "", Function: "Binds copper, nickel, and fatty acids as well as, … the human AFP shows estrogen-binding properties.", GeneName: "AFP", Kegg: "hsa:174;", …}
    7: {ChemblId: "CHEMBL3594;", EcNumber: "4.2.1.1", Function: "Reversible hydration of carbon dioxide. Participat…ervical neoplasia. {ECO:0000269|PubMed:18703501}.", GeneName: "CA9", Kegg: "hsa:768;", …}
    8: {ChemblId: "CHEMBL1764938;", EcNumber: "", Function: "Calcium-regulated membrane-binding protein whose a…ibits PCSK9-enhanced LDLR degradation, probably r", GeneName: "ANXA2", Kegg: "hsa:302;", …}
    9:
    ChemblId: ""
    EcNumber: ""
    Function: "Inhibitor of phospholipase A2, also possesses anti-coagulant properties. Also cleaves the cyclic bond of inositol 1,2-cyclic phosphate to form inositol 1-phosphate."
    GeneName: "ANXA3"
    Kegg: "hsa:306;"
    Length: 323
    Mass: 36375
    Name: "Annexin A3 "
    Pdb: "1AII;1AXN;"
    ProteinFamily: "Annexin family"
    UniprotId: "P12429"
```

问题:

3. MOA和Structures部分,全部和TCMAnalyzer一样么?


4. 数据方面:
    - 数据库中空值置Null 而不是 空字符串
    - Target部分的外链数据末尾多了分号
    - 每个target需要一个ID, 并绑定内部访问链接
    - Ingredients中 LigandId 字段是什么意思


```json
{
    "draw": 1,
    "recordsTotal" : 497,
    "recordsFiltered" : 100,
    "msg": "",
    "state": "success",
    "data": [
        {...}
    ]
}
```