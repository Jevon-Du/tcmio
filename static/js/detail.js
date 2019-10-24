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
        render(msg.Data, type);
        if ('Mol' in msg.Data){
            show_mol_structure(msg.Data.Mol);
        }
    });

});

function which_type(data_type){
    var name=capitalizeFirstLetter(data_type.slice(0,-1));
    return name;
}

function render(data, data_type){
    var name = ('Name' in data ? data.Name : data.Id);
    var heading = which_type(data_type) + ' : ' + name;
    $('.item_name').html(heading);

    // 创建表格,渲染数据
    var html = '';
    for (var key in data){
        html += '<tr class="item">';
        html += '<th class="item_name">' + key +'</th>';
        if (data[key]){
            if(key=='Mol'){
                html += '<td class="item_info">' + canvas_format(key) +'</td>';
            } else{
                html += '<td class="item_info">' + link_format(data[key], key) +'</td>';
            }
        } else {
            html += '<td class="item_info"><i>N/A<i></td>';
        }
        html += '</tr>';
    };

    $('.info tbody').append(html);
}

function show_mol_structure(molstr){
    var myCanvas = new ChemDoodle.ViewerCanvas('sketcherMol', 300, 300);
    myCanvas.emptyMessage = 'No Data Loaded!';
    myCanvas.repaint();
    var mol = ChemDoodle.readMOL(molstr);
    myCanvas.loadMolecule(mol);
}