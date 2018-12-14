$(function () {

    //1.初始化Table
    var oTable = new TableInit();
    oTable.Init();

    //2.初始化Button的点击事件
    var oButtonInit = new ButtonInit();
    oButtonInit.Init();

});


var TableInit = function () {
    var oTableInit = new Object();
    //初始化Table
    oTableInit.Init = function () {
        $('#user-manager-tb').bootstrapTable({
            url: '/mysql/getsqlusermanager',         //请求后台的URL（*）
            method: 'get',                      //请求方式（*）
            toolbar: '#sqluser-tool',                //工具按钮用哪个容器
            striped: true,                      //是否显示行间隔色
            cache: false,                       //是否使用缓存，默认为true，所以一般情况下需要设置一下这个属性（*）
            pagination: true,                   //是否显示分页（*）
            sortable: true,                     //是否启用排序
            sortOrder: "asc",
            sortName: "Name",
            //排序方式
            //queryParams: oTableInit.queryParams,//传递参数（*）
            /*
            分页方式：client客户端分页，server服务端分页（*）https://www.cnblogs.com/landeanfen/p/4976838.html,
            如果错误选择分页方式无法正常显示
            */
            sidePagination: "client",
            pageNumber:1,                       //初始化加载第一页，默认第一页
            pageSize: 10,                       //每页的记录行数（*）
            pageList: [10, 25, 50, 100],        //可供选择的每页的行数（*）
            search: true,                       //是否显示表格搜索，此搜索是客户端搜索，不会进服务端，所以，个人感觉意义不大
            strictSearch: true,
            showColumns: true,                  //是否显示所有的列
            showRefresh: true,                  //是否显示刷新按钮
            minimumCountColumns: 2,             //最少允许的列数
            clickToSelect: true,                //是否启用点击选中行
            height: 500,                        //行高，如果没有设置height属性，表格自动根据记录条数觉得表格高度
            uniqueId: "Name",                     //每一行的唯一标识，一般为主键列
            showToggle:true,                    //是否显示详细视图和列表视图的切换按钮
            cardView: false,                    //是否显示详细视图
            detailView: false,                   //是否显示父子表
            columns: [{
                checkbox: true
            }, {
                field: 'Name',
                title: '用户名'
            }, {
                field: 'Department',
                title: '部门'
            }, {
                field: 'Roleid',
                title: '角色'
            }, ]
        });
    };

    //得到查询的参数
    oTableInit.queryParams = function (params) {
        var temp = {   //这里的键的名字和控制器的变量名必须一直，这边改动，控制器也需要改成一样的
            limit: params.limit,   //页面大小
            offset: params.offset,  //页码
            departmentname: $("#txt_search_departmentname").val(),
            statu: $("#txt_search_statu").val()
        };
        return temp;
    };
    return oTableInit;
};


var ButtonInit = function () {
    var oInit = new Object();
    var postdata = {};

    oInit.Init = function () {
        //初始化页面上面的按钮事件
    };

    return oInit;
};

//删除用户
$("#btn_delete").click(function () {
     if ($('input[name="btSelectItem"]:checked ').length > 0 ){
         layer.confirm('确定删除用户？', {
             btn: ['删除','取消'] //按钮
         }, function(){
             $('input[name="btSelectItem"]:checked ').each(function () {
                 //alert($(this).parent().next().text())
                 var username = $(this).parent().next().text()
                 var table = "user"
                 var token = $.cookie('token')
                 $.ajax({
                     async: true,
                     type:'get',
                     url:'/mysql/delete?table=' + table + '&username=' + username,
                     dataType:'json',
                     beforeSend: function(request) {
                         request.setRequestHeader("Authorization", token);
                     },
                     success:function (data){
                         //抛出删除成功，并不显示此行
                         layer.msg("用户已删除",{
                             icon: 1,
                             time: 1000,
                             end:function () {
                                 //如何实现局部div刷新
                                 window.location.reload()
                             }
                         });
                         //  $(this).parent().remove();
                         //  $(this).css("display","none") 为什么不生效
                     },
                     error:function(){
                         alert('inner error')
                     }
                 })
             })
         });
         //var selectIndex = $('input[name="btSelectItem"]:checked ').parent().next().text();
     } else {
         console.info("donot select any el")
     }
})
//在线编辑用户
$("#btn_edit").click(function () {
    if ($('input[name="btSelectItem"]:checked ').length == 1 ){
        $('#editSqlUserbtn').modal("show");
    } else  {
        layer.msg("请您选择一位用户",{
            icon: 1,
            time: 2000,
        });
    }
})

//模态框打开在线编辑时，自动填入文档
$("#editSqlUserbtn").on('show.bs.modal',function () {
    if ($('input[name="btSelectItem"]:checked ').length == 1 ){
        var username = $('input[name="btSelectItem"]:checked ').parent().next().text();
        var department = $('input[name="btSelectItem"]:checked ').parent().next().next().text();
        var roleid = $('input[name="btSelectItem"]:checked ').parent().next().next().next().text();
        $("#sql-username").val(username)
        $("#sql-department").val(department)
        $("#sql-roleid").val(roleid)
    }
})
//模态框关闭时，提交按钮添加disabled属性
$("#editSqlUserbtn").on('hide.bs.modal',function () {
    $('#editSqlUser').attr("disabled","disabled");
})

$(".edit-user-ipt").change(function (){
    //$('#editSqlUser').data('changed',true);
    $('#editSqlUser').removeAttr("disabled");

});
//编辑用户
$("#editSqlUser").click(function () {
    var username =  $("#sql-username").val()    //唯一值
    var department = $("#sql-department").val()
    var roleid = $("#sql-roleid").val()
    $.ajax({
        async: true,
        type:'get',
        url:'/user/sql/edit?username=' + username + '&department=' + department + '&roleid=' + roleid,
        dataType:'json',
        success:function (data){
            //抛出删除成功，并不显示此行
            layer.msg("用户信息已更新",{
                icon: 1,
                time: 1000,
                end:function () {
                    //如何实现局部div刷新
                    window.location.reload()
                }
            });
            //  $(this).parent().remove();
            //  $(this).css("display","none") 为什么不生效
        },
        error:function(){
            alert('inner error')
        }
    })
})
//新增用户
$("#addSqlUser").click(function () {
    var username =  $("#add-sql-username").val()    //唯一值
    var password =  $("#add-sql-password").val()
    var department = $("#add-sql-department").val()
    var roleid = $("#add-sql-roleid").val()
    $.ajax({
        async: true,
        type:'get',
        url:'/user/sql/add?username=' + username + '&password=' +   password + '&department=' +   department + '&roleid=' + roleid,
        dataType:'json',
        success:function (data){
            //抛出删除成功，并不显示此行
            layer.msg("创建用户成功",{
                icon: 1,
                time: 500,
                end:function () {
                    //如何实现局部div刷新
                    window.location.reload()
                }
            });
            //  $(this).parent().remove();
            //  $(this).css("display","none") 为什么不生效
        },
        error:function(){
            alert('inner error')
        }
    })
})