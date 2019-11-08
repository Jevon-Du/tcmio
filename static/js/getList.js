$(document).ready( function () {
    var html_type = window.location.pathname.split('/')[2];
    console.log(html_type);
    $('.breadcrumb .active').html(capitalizeFirstLetter(html_type));

    // 初始化target表格
    initialize_table(html_type);
    $('.dataTable').show(100);

    // 设置class=active
    $("#navbar .dropdown-menu li").removeClass("active");
    $("#navbar .dropdown-menu ."+html_type).addClass("active");
});

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
        dom: '<"toolbar"l><r<t>ip>',
        ordering: true,
        scrollX: true,
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

   return table;
}