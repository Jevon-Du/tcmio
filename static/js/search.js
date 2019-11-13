/*
**  global variable
*/
//  存储所有的节点数组和边数组的对象
var analyzeResult=null;

//  存储当前选中节点的box位置，click时触发
var nodePostion=null;

//  存储当前选中节点时所展示的详细信息
var detailNodeInfo=null;

//  存储当前选中节点的Group
var selectedNodeGroup=null;

var DIVCOUNT = 0;
var RECORDS = 5;
var MOLS=[];

var TABLE1 = null;
var TABLE2 = null;
var TABLES = [];

// var COL_TYPE = {
//     targets: [
//         {"data": "Id", "className":"innerLink", "width":"6%"},
//         {"data": "Name"},
//         {"data": "GeneName"},
//         {"data": "Function", "visible": false},
//         {"data": "ProteinFamily"},
//         {"data": "UniprotId"},
//         {"data": "ChemblId"},
//         {"data": "EcNumber"}
//         // {"data": "Kegg"},
//         // {"data": "Pdb"},
//         // {"data": "Length"},
//         // {"data": "Mass"}
//     ],
//     ligands: [
//         {"data": "Id", "className":"innerLink", "width":"6%"},
//         {"data": "Mol", "width":"200px", "className":"chem-viewer"},
//         {"data": "Name"},
//         {"data": "IngredientId"},
//         {"data": "ChemblId"},
//         // {"data": "Smiles", "visible": false},
//         // {"data": "Inchi", "visible": false},
//         // {"data": "Inchikey", "visible": false},
        
//         {"data": "Formula", "visible": false},
//         {"data": "MolWeight", "visible": false},
//         {"data": "Hba", "visible": false},
//         {"data": "Hbd", "visible": false},
//         {"data": "Rtb", "visible": false},
//         {"data": "Alogp", "visible": false}     
//     ],
//     ingredients: [
//         {"data": "Id", "className":"innerLink", "width":"6%"},
//         {"data": "Mol", "width":"200px", "className":"chem-viewer"},
//         {"data": "Name"},
//         // {"data": "Synonyms", "visible": false},
//         {"data": "Smiles"},
//         // {"data": "Inchi", "visible": false},
//         // {"data": "Inchikey", "visible": false},      
//         {"data": "LigandId", "visible": false}
//     ],
//     tcms: [
//         {"data": "Id", "className":"innerLink", "width":"6%"},
//         {"data": "ChineseName"},
//         {"data": "PinyinName"},
//         {"data": "EnglishName"},
//         {"data": "UsePart"},
//         {"data": "PropertyFlavor"},
//         {"data": "ChannelTropism"},            
//         {"data": "Effect"},
//         {"data": "Indication"},
//         {"data": "RefSource"}
//     ],
//     prescriptions: [
//         {"data": "Id", "className":"innerLink", "width":"6%"},
//         {"data": "ChineseName"},
//         {"data": "PinyinName"},
//         {"data": "Ingredients"},
//         {"data": "Indication"},
//         {"data": "Effect"},
//         {"data": "RefSource"}
//     ]
// }


var selectedName = 'english_name';
var searchType = 'tcms';
var name_datas = {
    'tcmsenglish_name': [],
    'tcmspinyin_name': [],
    'tcmschinese_name': [],
    'prescriptionspinyin_name': [],
    'prescriptionschinese_name': []
};

$(document).ready( function () {
    var search_type = window.location.pathname.slice(1);
    if (search_type=='structure'){
        structure_search();
    } else if(search_type=='moa'){
        moa_search();
    }
});

function structure_search(){
    $('.methods ul a').on('click', function () {
        var method = $(this).text();
        $('.method').html(method);
        if (method=='Similarity'){
            $('.thresholds button').removeClass('disabled');
        } else {
            $('.thresholds button').addClass('disabled');
        }
    });
    
    $('.thresholds ul a').on('click', function () {
        $('.threshold').html($(this).text());
    });
    $('.types ul a').on('click', function () {
        $('.type').html($(this).text());
    });

    var method_map = {
        'Fullstructure': 'exact',
        'Similarity': 'sim',
        'Substructure': 'sub'
    };

    $('.search a').on({
        click: function(){
            var iframe = document.getElementById('iframe');
            var query = iframe.contentWindow.ipmDraw.getMol();
            var method = $('.methods .method').text();
            var threshold = $('.thresholds .threshold').text().slice(3);
            var type = $('.types .type').text().toLowerCase(); // ingredient or ligand

            var request = {
                "query": query,
                "method": method,
                "threshold": threshold
            }

            var data_type = type+'s';

            // 方法一: 如果表格已经初始化, 则destory, 但是这种方式涉及大量dom操作,影响性能
            if (TABLE1 & TABLE2){
                TABLE1.destroy(false); //false: 保留原始的DOM
                TABLE2.destroy(false);
                TABLE1 = null;
                TABLE2 = null;
            }

            // 初始化查询结果表格
            TABLE1 = initialize_table_in_browse(data_type, request, false, true, '<r<t>ip>');
            // 初始化分析结果表格
            TABLE2 = initialize_table_in_structure(data_type, request);
            
            if (data_type=='ingredients') {

                $('#ligands_wrapper').hide();
                $('#ligands_analyze_wrapper').hide();

                $('#ingredients').show(100);
                $('#ingredients_wrapper').show(100);
                $('#ingredients_analyze').show(100);
                $('#ingredients_analyze_wrapper').show(100);

            } else if (type=='ligand') {

                $('#ingredients_wrapper').hide();
                $('#ingredients_analyze_wrapper').hide();

                $('#ligands').show(100);
                $('#ligands_wrapper').show(100);
                $('#ligands_analyze').show(100); 
                $('#ligands_analyze_wrapper').show(100);           
            }

            $('#collapseOne').collapse('hide');
            $('#collapseTwo').collapse('show');
            $('#collapseThree').collapse('show');

        }
    });
}

function moa_search(){
    // setp1: 请求默认数据
    var data = get_names();
    name_datas[searchType+selectedName] = data;
    // step2: 显示select
    $selectEl = $('.'+searchType+selectedName);
    $selectEl.show();
    // setp3: 插入数据到option
    $selectEl.append(gene_option(data, selectedName));
    // step4: 初始化select
    initial_select2($selectEl);

    $('.searchTypes ul a').on('click', function () {
        var type = $(this).text();
        // 先清空以前选中的数据
        $('.'+searchType+selectedName).val(null).trigger('change');
        // 更新searchType
        searchType = $(this).attr('value');
        $('.searchType').html(type);
        if (searchType=='prescriptions'){
            $('.searchParam').html('PinyinName');
            selectedName = 'pinyin_name';
            $('.searchParams li.enName').addClass('disabled');
        } else {
            $('.searchParams li.enName').removeClass('disabled');
        }

        if (name_datas[searchType+selectedName].length == 0){
            // setp1: 请求默认数据
            name_datas[searchType+selectedName] = get_names();
            // step2: 显示select
            $('.searchValue select').hide();
            $('.searchValue .select2-container').hide();
            $selectEl = $('.'+searchType+selectedName);
            $selectEl.show();
            // setp3: 插入数据到option
            $selectEl.append(gene_option(name_datas[searchType+selectedName], selectedName));
            // step4: 初始化select
            initial_select2($selectEl);
        } else {
            // step2: 显示select
            $('.searchValue select').hide();
            $('.searchValue .select2-container').hide();
            $selectEl = $('.'+searchType+selectedName);
            $selectEl.show();
        }
    });
    
    $('.searchParams ul a').on('click', function () {
        if ($(this).parent().hasClass('disabled')){
            return false;
        }
        $('.searchParam').html($(this).text());
        // 先清空以前选中的数据
        $('.'+searchType+selectedName).val(null).trigger('change');
        // 更新selectedName        
        selectedName = $(this).attr('value');
        if (name_datas[searchType+selectedName].length == 0){
            // setp1: 请求默认数据
            name_datas[searchType+selectedName] = get_names();
            // step2: 显示select
            $('.searchValue select').hide();
            $('.searchValue .select2-container').hide();
            $selectEl = $('.'+searchType+selectedName);
            $selectEl.show();
            // setp3: 插入数据到option
            $selectEl.append(gene_option(name_datas[searchType+selectedName], selectedName));
            // step4: 初始化select
            initial_select2($selectEl);
        } else {
            // step2: 显示select
            $('.searchValue select').hide();
            $('.searchValue .select2-container').hide();
            $selectEl = $('.'+searchType+selectedName);
            $selectEl.show();
        }
    });

    // draw(null, null);

    $('.search a').on({
        click: function(){
            var search_val = [];
            $('.'+searchType+selectedName+' .select2-search-choice div').each(function(index,el){
                search_val.push($(el).html());
            });

            draw(network_data2['nodes'], network_data2['edges']);

            // $.ajax({
            //     url: 'network/' + searchType,
            //     type: "GET",
            //     dataType: "html",
            //     data:{
            //         'kw': search_val.join(','),
            //         'type': selectedName
            //     },
            //     error: function (jqXHR, textStatus, errorThrown) {
            //         if (textStatus == "timeout") {
            //             alert("Request timeout, please refresh the page.");
            //         } else {
            //             alert(textStatus);
            //         }
            //     }
            // }).done(function (msg) {
            //     console.log(msg);
            //     // TODO: 渲染数据
            //     draw(msg.Data['nodes'], msg.Data['edges']);
            // });

            $('#collapseOne').collapse('hide');
            $('#collapseTwo').collapse('show');
            

        }
    });
}

function get_names(){
    var data=null;
    $.ajax({
        url: '/'+searchType,
        type: "get",
        dataType: "json",
        async: false,
        error: function (jqXHR, textStatus, errorThrown) {
            if (textStatus == "timeout") {
                alert("Request timeout, please refresh the page.");
            } else {
                alert(textStatus);
            }
        }
    }).done(function (msg) {
        data = msg.data;
    });
    return data;
}

function initial_select2($el){
    $($el).select2({
        // theme: "flat",
        dropdownCssClass: 'dropdown-inverse',
        width: '100%',
        allowClear: true,
        placeholder: 'Select a name or some names',
    });
}


function gene_option(data, type){
    var html='';
    var name_map = {
        'english_name': 'EnglishName',
        'pinyin_name': 'PinyinName',
        'chinese_name': 'ChineseName'
    }
    var _type = name_map[type];
    data.forEach(function(currentValue, index){
        html = html + '<option value="' + data[index]['Id'] + '">'+ data[index][_type] +'</option>';
    });
    return html;
}

function initialize_table(data_type){
    var $table = $('#' + data_type);

    var table = $table.DataTable({
        serverSide: true,
        autoWidth: true,
        ajax: {
            url: '/'+data_type,
            type: 'get',
            dataSrc: function (msg) {
                return data_format(msg, data_type);
            }
        },
        columns: COL_TYPE[data_type],
        dom: '<r<t>ip>',
        ordering: true,
        scrollX: true,
        pagingType:   "full_numbers",
        pageLength: 5, //每页显示的初始记录数量
        language: {
            "lengthMenu": "每页 _MENU_ 条记录",
            "zeroRecords": "没有找到记录",
            "info": "第 _PAGE_ 页 ( 总共 _PAGES_ 页 )",
            "infoEmpty": "无记录",
            "infoFiltered": "(从 _MAX_ 条记录过滤)"
        }
    });
    // 翻页事件的处理
    $table.on( 'page.dt', function () {
        // var info = table.page.info();
        // $('#pageInfo').html( 'Showing page: '+info.page+' of '+info.pages );
        init_chemdoodle(data_type);
    });

    // 排序事件处理
    $table.on( 'order.dt', function () {
        // This will show: "Ordering on column 1 (asc)", for example
        // var order = table.order();
        // $('#orderInfo').html( 'Ordering on column '+order[0][0]+' ('+order[0][1]+')' );
        init_chemdoodle(data_type);
    });

   init_chemdoodle(data_type);
}
/**************************** Data fromat ****************************/