//checker page use this jquery script
//sleep 10 second to get urged msg use ajax
    //审核通过
    $(".btn-check-ok").click(function () {

        var  sqlid =$(this).parent().siblings(".sql-id").text()
        sqlid = sqlid.replace(/\s/g, "") ;
        $.ajax({
            type: "POST",
            url: "/mysql/update",
            data: {
                sqlid: sqlid,
                checkstatus: 1,
            },
            dataType:'json',
            error: function(data){
                alert(data.Status);
            },
            success:function(){
                layer.msg("审核成功",{
                    icon: 1,
                    time: 1000,
                    end:function () {
                        //如何实现局部div刷新
                        //window.location.reload()
                        console.info("urged other do that ")

                    }
                });
            }
        });
    })
    //驳回
    $(".btn-check-false").click(function () {

    var  sqlid =$(this).parent().siblings(".sql-id").text()
    sqlid = sqlid.replace(/\s/g, "") ;
    $.ajax({
        type: "POST",
        url: "/mysql/update",
        data: {
            sqlid: sqlid,
            checkstatus: 2,
        },
        dataType:'json',
        error: function(data){
            alert(data.Status);
        },
        success:function(){
            layer.msg("任务已驳回",{
                icon: 1,
                time: 1000,
                end:function () {
                    //如何实现局部div刷新
                    //window.location.reload()
                    console.info("urged other do that ")

                }
            });
        }
    });
})
