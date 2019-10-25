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
    }
    $('.search a').on({
        click: function(){
            var iframe = document.getElementById('iframe');
            var query = iframe.contentWindow.ipmDraw.getMol();
            var method = $('.methods .method').text();
            var threshold = $('.thresholds .threshold').text().slice(3);
            var type = $('.types .type').text().toLowerCase();

            $.ajax({
                url: 'structure/'+type,
                type: "GET",
                dataType: "html",
                data:{
                    'query': query,
                    'method': method_map[method],
                    'threshold': threshold
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
                // TODO: 渲染数据
            });
            $('#collapseOne').collapse('hide');
            $('#collapseTwo').collapse('show');
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

    draw(null, null);

    $('.search a').on({
        click: function(){
            var search_val = [];
            $('.'+searchType+selectedName+' .select2-search-choice div').each(function(index,el){
                search_val.push($(el).html());
            });
            console.log(search_val.join(','));
            console.log(selectedName);
            console.log(searchType);

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
                console.log(msg);
                // TODO: 渲染数据
            });

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
