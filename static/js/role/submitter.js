$(document).ready(function(){
    //sql审批任务撤回
    $(".row-delete").click(function(){
        //后续需要传入DBname ,
       var  sqlid =$(this).parent().siblings(".sql-id").text()
        sqlid = sqlid.replace(/\s/g, "") ;
       var table = "submit_status"
       var token = $.cookie('token')
        $.ajax({
            async: true,
            type:'get',
            url:'/mysql/delete?id=' + sqlid + '&table=' + table,
            dataType:'json',
            beforeSend: function(request) {
                request.setRequestHeader("Authorization", token);
            },
            success:function (data){
                //抛出删除成功，并不显示此行
                layer.msg("任务已撤销",{
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
            });
    //sql审批任务催促
    $(".row-urged").click(function(){
        //后续需要传入DBname ,
        var  sqlid =$(this).parent().siblings(".sql-id").text()
        sqlid = sqlid.replace(/\s/g, "") ;
        var pusher = $(this).parent().siblings(".Approver").text()
        //执行人是否取自表格，需要定义
        var reciver = "qiushu"
        $.ajax({
            type: "POST",
            url: "/mysql/urged",
            data: {
                pusher: pusher,
                sqlid: sqlid,
                reciver: reciver,
            },
            dataType:'json',
            error: function(data){
                alert(data.Status);
            },
            success:function(data){
                layer.msg("已催促对方审批",{
                    icon: 1,
                    time: 1000,
                    end:function () {
                        //如何实现局部div刷新
                        window.location.reload()


                    }
                });
            }
        });

    });

    //sql审批提示状态
    /*
    $(".checker-tip").mouseover(function () {
        layer.tips('0：已退回 1：已审批',".checker-tip",{
            tips: [1,'#286090'],
            tipsMore: false,
            time: 1000
        });
    });
    //sql执行状态
    $(".ops-tip").mouseover(function () {
        layer.tips('0：已退回 1：已执行',".ops-tip",{
            tips: [1,'#f39c12'],
            tipsMore: false,
            time: 1000
        });
    });
    */
    //sql 输入框textarea 高度自适应
    $('textarea').each(function () {
        this.setAttribute('style', 'height:' + '50' + 'px;overflow-y:hidden;');  //设置默认textarea行高
    }).on('input', function () {
        this.style.height = 'auto';
        this.style.height = (this.scrollHeight) + 'px';
    });
    //ajax 提交sql任务请求
    $("#addSqlBtn").click(function () {
        var dbname = $("#dbName").val()
        var checkername = $("#checkerName").val()
        var sqlstatement = $("#sqlStatement").val()
        var token = $.cookie('token')
        //var url = '/mysql/insert' + '?dbname=' + dbname.replace(/\s/g, "") + '&checkername=' + checkername.replace(/\s/g, "") + '&sqlstatement=' + sqlstatement.replace(/\s/g, "")
        $.ajax({
            type: "POST",
            url: "/mysql/insert",
            data: {
                dbname: dbname,
                checkername: checkername,
                sqlstatement: sqlstatement,
            },
            dataType:'json',
            beforeSend: function(request) {
                request.setRequestHeader("Authorization", token);
            },
            error: function(data){
                console.info("here")
                alert(data.Status);
            },
            success:function(data){
                layer.msg("提交成功",{
                    icon: 1,
                    time: 1000,
                    end:function () {
                        //如何实现局部div刷新
                        window.location.reload()
                    }
                });
            }
        });

    });
    });
