$(document).ready(function(){
    $(".row-ops-refuse").click(function () {

        var  sqlid =$(this).parent().siblings(".sql-id").text()
        sqlid = sqlid.replace(/\s/g, "") ;
        $.ajax({
            type: "POST",
            url: "/mysql/update",
            data: {
                sqlid: sqlid,
                execstatus: 2,
            },
            dataType:'json',
            error: function(data){
                alert(data.Status);
            },
            success:function(){
                layer.msg("任务已拒绝",{
                    icon: 1,
                    time: 1000,
                    end:function () {
                        //如何实现局部div刷新
                        //window.location.reload()
                        console.info("任务已被管理员拒绝 ")
                    }
                });
            }
        });
    })
    });
