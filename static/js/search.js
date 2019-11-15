/*
**  global variable
*/

var TABLE1 = null;
var TABLE2 = null;
var TABLES = [];

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
                "method": method_map[method],
                "threshold": threshold
            }

            console.log(request);

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
                $('#result .panel-title a').html('Hits');
                $('#activity .panel-title a').html('TCM');

                $('#ligands_wrapper').hide();
                $('#ligands_analyze_wrapper').hide();

                $('#ingredients').show(100);
                $('#ingredients_wrapper').show(100);
                $('#ingredients_analyze').show(100);
                $('#ingredients_analyze_wrapper').show(100);

            } else if (type=='ligand') {
                $('#result .panel-title a').html('Ligand Hits');
                $('#activity .panel-title a').html('Activity');

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

    $('.search a').on({
        click: function(){
            var search_val = [];
            $('.'+searchType+selectedName+' .select2-search-choice div').each(function(index,el){
                search_val.push($(el).html());
            });

            // 加载本地json数据绘制网络
            // draw(network_data2['nodes'], network_data2['edges']);
            // $('#collapseOne').collapse('hide');
            // $('#collapseTwo').collapse('show');
            
            // 请求网络数据
            $.ajax({
                url: 'network/' + searchType,
                type: "GET",
                dataType: "html",
                data:{
                    'kw': search_val.join(','),
                    'type': selectedName
                },
                error: function (jqXHR, textStatus, errorThrown) {
                    if (textStatus == "timeout") {
                        alert("Request timeout, please refresh the page.");
                    } else {
                        alert(textStatus);
                    }
                }
            }).done(function (msg) {
                var msg = JSON.parse(msg);
                console.log(msg);
                draw(msg.Data['nodes'], msg.Data['edges']);

                $('#collapseOne').collapse('hide');
                $('#collapseTwo').collapse('show');
            
            });
        }
    });

    //点击pathway请求数据
    $('#collapseThree').on({
        click: function(){
            var search_val = [];
            $('.'+searchType+selectedName+' .select2-search-choice div').each(function(index,el){
                search_val.push($(el).html());
            });
            // 请求pathway数据
            $.ajax({
                url: '/',
                type: "GET",
                dataType: "html",
                data:{
                    'kw': search_val.join(','),
                    'type': selectedName
                },
                error: function (jqXHR, textStatus, errorThrown) {
                    if (textStatus == "timeout") {
                        alert("Request timeout, please refresh the page.");
                    } else {
                        alert(textStatus);
                    }
                }
            }).done(function (msg) {
                console.log(msg);
                var $table_pathways = $('#pathways').DataTable();
                initialize_table_in_pathway($table_pathways, msg.Data);
                $('#pathways').show(100);
                $('#pathways').show(100);
                // $('#collapseThree').collapse('show');
            });
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

