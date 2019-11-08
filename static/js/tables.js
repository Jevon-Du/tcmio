/*!
 * 先导入commmon.js: data_format, init_chemdoodle
 */

var DIVCOUNT = 0;
var RECORDS = 5;
var MOLS=[];

var COL_TYPE = {
    targets: [
        {"data": "Id", "className":"innerLink", "width":"10%"},
        {"data": "Name"},
        {"data": "GeneName"},
        {"data": "Function", "visible": false},
        {"data": "ProteinFamily"},
        {"data": "UniprotId"},
        {"data": "ChemblId"},
        {"data": "EcNumber"}
        // {"data": "Kegg"},
        // {"data": "Pdb"},
        // {"data": "Length"},
        // {"data": "Mass"}
    ], 
    ligands: [
        {"data": "Id", "className":"innerLink", "width":"10%"},
        {"data": "Mol", "width":"200px", "className":"chem-viewer"},
        {"data": "Name"},
        {"data": "IngredientId"},
        {"data": "ChemblId"},
        // {"data": "Smiles", "visible": false},
        // {"data": "Inchi", "visible": false},
        // {"data": "Inchikey", "visible": false},
        {"data": "Formula", "visible": false},
        {"data": "MolWeight", "visible": false},
        {"data": "Hba", "visible": false},
        {"data": "Hbd", "visible": false},
        {"data": "Rtb", "visible": false},
        {"data": "Alogp", "visible": false}     
    ],
    ingredients: [
        {"data": "Id", "className":"innerLink", "width":"10%"},
        {"data": "Mol", "width":"200px", "className":"chem-viewer"},
        {"data": "Name"},
        // {"data": "Synonyms", "visible": false},
        {"data": "Smiles"},
        // {"data": "Inchi", "visible": false},
        // {"data": "Inchikey", "visible": false},      
        {"data": "LigandId", "visible": false}
    ],
    tcms: [
        {"data": "Id", "className":"innerLink", "width":"6%"},
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
        {"data": "Id", "className":"innerLink", "width":"10%"},
        {"data": "ChineseName"},
        {"data": "PinyinName"},
        {"data": "Ingredients"},
        {"data": "Indication"},
        {"data": "Effect"},
        {"data": "RefSource"}
    ],
    ligands_analyze: [
        {"data": "Id", "className":"innerLink", "width":"10%"},
        {"data": "TargetChemblId", "className":"tChemblLink"},
        {"data": "MolChemblId", "className":"cChemblLink"},
        {"data": "Type"},
        {"data": "Value"},
        {"data": "Unit"},
        {"data": "RefId"}
    ],
    ingredients_analyze: [
        {"data": "ChineseName"},
        {"data": "IngredientID", "className":"innerLink"},
        {"data": "Count"}
    ]
};

function initialize_table_in_browse(data_type, request, has_button, has_mol, dom_){
    /**
     * 动态请求数据, 使用data_type构造url
     * @param data_type: str, 标示是哪类数据, e.g : ligands, targets, ingredients...
     * @param request: object, 额外的发送到服务器的参数
     * @param has_button: boolean, 标示table中是否有button功能区, default true
     * @param has_mol: boolean, 标示table中是否渲染分子结构, default false
     * @param dom_: str, toolbar dom配置, default '<"toolbar"l><r<t>ip>'
     * @return {table} DataTable对象
     */

    // data_type: 被初始化的table_id
    var $table = $('#' + data_type);

    var table = $table.DataTable({
        serverSide: true,
        autoWidth: true,
        ajax: {
            url: '/'+data_type,
            type: 'get',
            data: function ( d ) {
                return $.extend( {}, d, request);
            },
            dataSrc: function (msg) {
                return data_format(msg, data_type);
            }
        },
        columns: COL_TYPE[data_type],
        dom: dom_,
        ordering: true,
        pagingType: "full_numbers",
        pageLength: 5, //每页显示的初始记录数量
        language: {
            "lengthMenu": "每页 _MENU_ 条记录",
            "zeroRecords": "没有找到记录",
            "info": "第 _PAGE_ 页 ( 总共 _PAGES_ 页 )",
            "infoEmpty": "无记录",
            "infoFiltered": "(从 _MAX_ 条记录过滤)"
        },
        // stateSave: true
    });

    // button功能初始化
    if (has_button){
        new $.fn.dataTable.Buttons( table, {
            buttons: [{
                    name: 'colvis',
                    extend: 'colvis',
                    className: 'btncolvis btn btn-inverse'
                }]
        });
        table.buttons( 0, '.btncolvis' ).containers().appendTo('.toolbar');

        new $.fn.dataTable.Buttons( table, {
            buttons: [{
                    name: 'pdf',
                    extend: 'pdf',
                    // text: 'Download as csv'
                    className: 'btnpdf btn btn-inverse'
                }]
        });
        table.buttons( 1, '.btnpdf' ).containers().appendTo('.toolbar');

        new $.fn.dataTable.Buttons( table, {
            buttons: [{
                    name: 'primary',
                    extend: 'csv',
                    // text: 'Download as csv'
                    className: 'btncsv btn btn-inverse'
                }]
        });
        table.buttons(2, '.btncsv' ).containers().appendTo('.toolbar');
        

        // show_hide column
        $('a.toggle-vis').on( 'click', function (e) {
            e.preventDefault();
    
            // Get the column API object
            var column = table.column( $(this).attr('data-column') );
    
            // Toggle the visibility
            column.visible(!column.visible());
        });

        // 页面长度事件的处理
        $table.on( 'length.dt', function ( e, settings, len ) {
            console.log( 'New page length: '+len );
            RECORDS = len;
            init_chemdoodle(data_type);
        });

    }
    
    // 有structure则动态渲染分子结构
    if (has_mol){
        // 翻页事件的处理
        $table.on( 'page.dt', function () {
            // var info = table.page.info();
            // $('#pageInfo').html( 'Showing page: '+info.page+' of '+info.pages );
            init_chemdoodle(data_type);
        } );

        // 排序事件处理
        $table.on( 'order.dt', function () {
            // This will show: "Ordering on column 1 (asc)", for example
            // var order = table.order();
            // $('#orderInfo').html( 'Ordering on column '+order[0][0]+' ('+order[0][1]+')' );
            init_chemdoodle(data_type);
        });

        init_chemdoodle(data_type);
    }
    return table;
};

function initialize_table_in_structure(data_type, request){
    /**
     * 动态请求数据, 使用data_type构造url
     * @param data_type: str, 标示是哪类数据, e.g : ligands, ingredients...
     * @param request: object, 额外的发送到服务器的参数
     * @return {table} DataTable对象
     */

    // data_type: ligands_analyze, or ingredients_analyze
    var $table = $('#' + data_type +'_analyze');

    var table = $table.DataTable({
        sserverSide: true,
        autoWidth: true,
        ajax: {
            url: '/structure/analyze/'+data_type,
            type: 'get',
            data: function ( d ) {
                return $.extend( {}, d, request);
            },
            dataSrc: function (msg) {
                return analyze_data_format(msg, data_type);
            }
        },
        // data: data_arrays,
        columns: COL_TYPE[data_type],
        dom: '<r<t>ip>',
        ordering: true,
        pagingType:   "full_numbers",
        pageLength: 6, //每页显示的初始记录数量
        lengthChange: false, //允许修改每页的记录数量
        language: {
            "lengthMenu": "每页 _MENU_ 条记录",
            "zeroRecords": "没有找到记录",
            "info": "第 _PAGE_ 页 ( 总共 _PAGES_ 页 )",
            "infoEmpty": "无记录",
            "infoFiltered": "(从 _MAX_ 条记录过滤)"
        }
    });

    // 采用本地json数据测试通过

    // var data_arrays = null;

    // if (data_type=='ingredients') {
    //     data_arrays = ingredient_msg.Data;
    // } else {
    //     data_arrays = ligand_msg.Data;
    // }
     
    // var table = $table.DataTable({
    //     sserverSide: false,
    //     autoWidth: true,
    //     data: data_arrays,
    //     columns: COL_TYPE[data_type+'_analyze'],
    //     dom: '<r<t>ip>',
    //     ordering: true,
    //     pagingType:   "full_numbers",
    //     pageLength: 6, //每页显示的初始记录数量
    //     lengthChange: false, //允许修改每页的记录数量
    //     language: {
    //         "lengthMenu": "每页 _MENU_ 条记录",
    //         "zeroRecords": "没有找到记录",
    //         "info": "第 _PAGE_ 页 ( 总共 _PAGES_ 页 )",
    //         "infoEmpty": "无记录",
    //         "infoFiltered": "(从 _MAX_ 条记录过滤)"
    //     }
    // });

    return table;
}

