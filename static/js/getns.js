$(document).ready(function(){
    //只执行一次，完美解决重复调用AJAX问题
    $("#get-k8s-ns").one("click",function () {

        $.ajax({
            async: true,
            type:'get',
            url:'/api/k8s/namespaces' ,
            dataType:'json',
            success:function (data){
                var thishtml = ''
                $.each(data,function(i,item){
                    thishtml = thishtml + '<li><a href="/home/k8s/namespace/' + item +  '"><i class="fa fa-circle-o"></i>' + item + '</a></li>'
                });
                $("#k8s-ns-show").html(thishtml)
            },
            error:function(){
                alert('inner error')
            }
        })

    });
});
