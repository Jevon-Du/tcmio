var DIVCOUNT = 0;
var RECORDS = 5;
var COL_TYPE = {
    targets: [
        {"data": "Name"},
        {"data": "GeneName"},
        {"data": "Function"},
        {"data": "ProteinFamily"},
        {"data": "UniprotId"},
        {"data": "ChemblId"},
        {"data": "EcNumber"},            
        {"data": "Kegg"},
        {"data": "Pdb"},
        {"data": "Length"},
        {"data": "Mass"}
    ],
    ingredients: [
        {"data": "Id"},
        {"data": "Name"},
        {"data": "Synonyms"},
        {"data": "Mol"},
        {"data": "Smiles"},
        {"data": "Inchi"},
        {"data": "Inchikey"},            
        {"data": "LigandId"}
    ],
    tcms: [
        {"data": "Id"},
        {"data": "ChineseName"},
        {"data": "PinyinName"},
        {"data": "EnglishName"},
        {"data": "UsePart"},
        {"data": "PropertyFlavor"},
        {"data": "ChannelTropism"},            
        {"data": "Effect"},
        {"data": "Indication"},
        {"data": "RefSource"}
    ],
    prescriptions: [
        {"data": "Id"},
        {"data": "ChineseName"},
        {"data": "PinyinName"},
        {"data": "Ingredients"},
        {"data": "Indication"},
        {"data": "Effect"},
        {"data": "RefSource"}
    ]
}
$(document).ready( function () {
    // 顶部导航按钮点击事件
    $("#navbar .dropdown-menu li a").on({
        click: function () {
            $("#navbar .dropdown-menu li").removeClass("active");

            var $liEl = $(this).parent();
            $liEl.addClass("active");
            var table_type = $(this).text();
            // 销毁所有DataTable实例并隐藏已初始化的
            var $tables = $('.container table');
            $tables.each(function(){
                $(this).DataTable().destroy(false);
                $(this).hide();
            });

            var $cur_table = $('#'+table_type);
            $cur_table.show(100, initialize_table);
        }
    });

    // 初始化target表格
    $("#navbar .dropdown-menu .active a").click();

    //给按钮绑定点击事件
    // $("#table_id_example_button").click(function () {
    //     var url = 'targets/id/' + '';
    //     var column1 = table.row('.selected').data().column1;
    //     var column2 = table.row('.selected').data().column2;
    //     alert("第一列内容："+column1 + "；第二列内容： " + column2);
    // });
});

function initialize_table(){
    var table_type = $('#navbar .dropdown-menu .active').text();
    var $table = $('#' + table_type);
    $('.breadcrumb .active').html(table_type);

    var cur_url = $('#navbar .dropdown-menu .active a').attr('href_');

    var table = $table.DataTable({
        serverSide: true,
        processing: true,
        destroy: true,
        ajax: {
            url: cur_url,
            type: 'get',
            // dataSrc: 接受到的服务器数据
            dataSrc: function (msg) {
                if (msg.state != "success"){
                    return null;
                }
                return data_format(msg, cur_url);
            }
        },
        columns: COL_TYPE[cur_url],
        // dom: '<"top"l><"toolbar">rt<"bottom"ip><"clear">',
        dom: '<"toolbar"l><r<t>ip>',
        // buttons: [{
        //     extend: 'columnsToggle',
        //     columns: '.e'
        // } ],
        ordering: true,
        scrollX: 1020,
        pagingType:   "full_numbers",
        pageLength: 5, //每页显示的初始记录数量
        lengthChange: true, //允许修改每页的记录数量
        lengthMenu: [ 5, 10, 25, 50], //每页可以显示的记录数量
        language: {
            "lengthMenu": "每页 _MENU_ 条记录",
            "zeroRecords": "没有找到记录",
            "info": "第 _PAGE_ 页 ( 总共 _PAGES_ 页 )",
            "infoEmpty": "无记录",
            "infoFiltered": "(从 _MAX_ 条记录过滤)"
        },
        stateSave: true
    });
    new $.fn.dataTable.Buttons( table, {
        buttons: [{
                name: 'pdf',
                extend: 'pdf',
                // text: 'Download as csv'
                className: 'btnpdf'
            }]
    });
    table.buttons( 0, '.btnpdf' ).containers().appendTo('.toolbar');
    new $.fn.dataTable.Buttons( table, {
        buttons: [{
                name: 'primary',
                extend: 'csv',
                // text: 'Download as csv'
                className: 'btncsv'
            }]
    });
    table.buttons( 1, '.btncsv' ).containers().appendTo('.toolbar');
    new $.fn.dataTable.Buttons( table, {
        buttons: [{
                name: 'colvis',
                extend: 'colvis',
                className: 'btncolvis'
            }]
    });
    table.buttons( 2, '.btncolvis' ).containers().appendTo('.toolbar');

    // 点击查看详情
    $table.on('click.dt',function() {
         // Get the column API object
         var column = table.column( $(this).attr('data-column') );
         $.ajax({
            url: "ligands/10",
            type: "get",
            dataType: "html",
            error: function (jqXHR, textStatus, errorThrown) {
                if (textStatus == "timeout") {
                    alert("Request timeout, please refresh the page.");
                } else {
                    alert(textStatus);
                }
            }
         }).done(function (response) {
            var msg = JSON.parse(response);
            console.log(msg);
        });
    });

    // show_hide column
    $('a.toggle-vis').on( 'click', function (e) {
        e.preventDefault();
 
        // Get the column API object
        var column = table.column( $(this).attr('data-column') );
 
        // Toggle the visibility
        column.visible(!column.visible());
    });
}
/**************************** Data fromat ****************************/

function data_format(msg, type){
    var col_name = [];
    if (type=='targets'){
        col_name = ['Name', 'GeneName', 'Function', 'ProteinFamily', 'UniprotId', 'ChemblId', 'EcNumber',
        'Kegg', 'Pdb', 'Mass', 'Length'];
        msg.data.forEach(function(currentValue, index){
            for (var idx in col_name) {
                var new_val = link_format(currentValue[col_name[idx]], col_name[idx]);
                msg.data[index][col_name[idx]] = new_val;
            }
        });
    } else {
        if(type=='ingredients'){
            col_name = ['Id', 'Name', 'Synonyms', 'Mol', 'Smiles', 'Inchi', 'Inchikey', 'LigandId'];
        } else if(type=='tcms'){
            col_name = ['Id', 'ChineseName', 'PinyinName', 'EnglishName', 'UsePart', 'PropertyFlavor', 'ChannelTropism',
            'Effect', 'Indication', 'RefSource'];
        } else if(type=='prescriptions'){
            col_name = ['Id', 'ChineseName', 'PinyinName', 'Ingredients', 'Indication', 'Effect', 'RefSource'];
        }
        msg.data.forEach(function(currentValue, index){
            for (var idx in col_name) {
                var new_val = link_format(currentValue[col_name[idx]], col_name[idx]);
                msg.data[index][col_name[idx]] = new_val;
            }
        });
        // msg.data.forEach(function(currentValue){
        //     var temp_a = [];
        //     for (var index in col_name) {
        //         temp_a.push(currentValue[col_name[index]]);
        //     };
        //     data_a.push(temp_a);
        // });
    }
    return msg.data;
};

function link_format(id_val, type){
    var db_link = {
        'UniprotId': 'https://www.uniprot.org/uniprot/',
        'ChemblId': 'https://www.ebi.ac.uk/chembl/target_report_card/',
        'EcNumber': 'https://enzyme.expasy.org/EC/',
        'Kegg': 'https://www.genome.jp/dbget-bin/www_bget?',
        'Pdb': 'https://www.rcsb.org/structure/'
    }
    var link_icon_el = '&nbsp<span class="glyphicon glyphicon-new-window" aria-hidden="true"></span>';
    if (db_link.hasOwnProperty(type)){
        var new_el = '';
        if (type=='UniprotId'){
            new_el = '<a href="https://www.uniprot.org/uniprot/'+ id_val +'" target="_blank">' + id_val + link_icon_el +'</a>';
        } else {
            var ids = id_val.split(';');
            for (var i in ids){
                new_el = new_el + '<a href="' + db_link[type] + ids[i].replace(' ', '') +'" target="_blank">' + ids[i].replace(' ', '') + link_icon_el +'</a>&nbsp';
            }
        }
        return new_el;
    }
    return id_val;
}