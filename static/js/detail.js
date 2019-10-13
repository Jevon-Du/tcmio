$(document).ready( function () {
    var html_type = window.location.pathname;

    // 请求数据
    $.ajax({
        url: html_type+'/json',
        type: "get",
        dataType: "json",
        error: function (jqXHR, textStatus, errorThrown) {
            if (textStatus == "timeout") {
                alert("Request timeout, please refresh the page.");
            } else {
                alert(textStatus);
            }
        }
    }).done(function (msg) {
        console.log(msg);
        var type= html_type.split('/')[1];

        render(msg.Data, which_type(type));
    });

});

function which_type(data_type){
    var name='';
    if (data_type=='targets'){
        name='Target';
    } else if(data_type=='ligands'){
        name='Ligand';
    } else if(data_type=='ingredients'){
        name='Ingredient';
    } else if(data_type=='tcms'){
        name='TCM';
    } else if(data_type=='prescriptions'){
        name='Prescription';
    }
    return name;
}

function render(data, data_type){
    var heading = data_type + ' : ' + data.Name;
    $('.item_name').html(heading);
    var html = '';
    // 创建表格,渲染数据
    for (var key in data){
        html += '<tr class="item">';
        html += '<th class="item_name">' + key +'</th>';
        if (data[key]){
            html += '<td class="item_info">' + link_format(data[key], key) +'</td>';
        } else {
            html += '<td class="item_info"><i>N/A<i></td>';
        }
        html += '</tr>';
    };

    $('.info tbody').append(html);
}

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
                new_el = new_el + '<a href="' + db_link[type] + ids[i].replace(' ', '') +'" target="_blank">' + ids[i].replace(' ', '') + link_icon_el +'</a>&nbsp&nbsp&nbsp';
            }
        }
        return new_el;
    }
    return id_val;
}
