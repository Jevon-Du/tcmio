/**********************format every link*******************/
function link_format(id_val, type){
    var db_link = {
        'UniprotId': 'https://www.uniprot.org/uniprot/',
        'ChemblId': 'https://www.ebi.ac.uk/chembl/target_report_card/',
        'EcNumber': 'https://enzyme.expasy.org/EC/',
        'Kegg': 'https://www.genome.jp/dbget-bin/www_bget?',
        'Pdb': 'https://www.rcsb.org/structure/',
        'Drug': 'https://www.drugbank.ca/drugs/'
    }
    var link_icon_el = '&nbsp<span class="glyphicon glyphicon-new-window" aria-hidden="true"></span>';
    if (db_link.hasOwnProperty(type)){
        var new_el = '';
        if (type=='UniprotId'){
            new_el = '<a href="https://www.uniprot.org/uniprot/'+ id_val +'" target="_blank">' + id_val + link_icon_el +'</a>';
        } else {
            var ids = id_val.split(';');
            for (var i in ids){
                new_el = new_el + '<a href="' + db_link[type] + ids[i].replace(' ', '') +'" target="_blank">' + ids[i].replace(' ', '') + link_icon_el +'</a>&nbsp&nbsp';
            }
        }
        return new_el;
    }
    return id_val;
}

function capitalizeFirstLetter(string) {
    return string.charAt(0).toUpperCase() + string.slice(1);
}