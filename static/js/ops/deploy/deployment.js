$(function () {
    //DOM加载后，初始化
    //初始化完毕

    $("#addsqljob-checkbox").click(function () {
        if ($(this).prop("checked")) {//jquery 1.6以前版本 用  $(this).attr("checked")
            console.info("print add ")
            $("#deploy-db-select").css("display","block")
        } else {
            console.info("print remove ")
            // $("#deploy-db-select").remove()
            $("[placeholder=\"input sql\"]").val("")

            $("#deploy-db-select").css("display","none")
        }
    });

    $("#adddeployjob-btn").click(function () {
        var onlineTitle = $("#online-title").val();
        var onlineContent = $("#online-content").val();
        var onlineItems = $("#online-items").val().join(",");
        var onlineDbname = $("#online-dbname").val();
        var onlineSqlStatement = $("#online-sqlstatement").val();
        console.info(onlineSqlStatement)
        $.ajax({
            type: "POST",
            url: "/api/deploy/additem",
            data: {
                onlineTitle: onlineTitle,
                onlineContent: onlineContent,
                onlineItems: onlineItems,
                onlineDbname: onlineDbname,
                onlineSqlStatement: onlineSqlStatement,
            },
            dataType:'json',
            //为什么会报inner错误
            error:function(data){
                console.info("inner error");
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
    
    $(".deploy-job-item").click(function () {
        $(this).siblings().removeClass("deploy-job-item-yellow");
        $(this).addClass("deploy-job-item-yellow");
        var lastjobtime = $(this).find("span").text()
        $.ajax({
                async: true,
                type:'get',
                url:'/api/deploy/getitemplan?lastjobtime=' + lastjobtime ,
                dataType:'json',
                success:function (data){
                    console.info("debug 01");

                    var jobItemResult = data;
                    var thishtml = ' <div class="box-header with-border">' + "\n" +
                        '<h3 class="box-title">' + jobItemResult.Title + ' </h3>' + "\n" +
                        '</div>' + "\n" +
                        '<div class="box-body" >' + "\n" +
                        '<strong><i class="fa fa-book margin-r-5"></i> 上线内容</strong>'   + "\n" +
                        '<p class="text-muted">'  + "\n" +
                        data.Content + "\n" +
                        '</p>' + "\n" +
                        '<hr>' + "\n" +
                        '<strong><i class="fa fa-sort-amount-asc margin-r-5"></i> 上线步骤</strong>' + "\n" +
                        '<ol>'  + "\n" ;
                    var step = data.DeploymentStep
                    for (var i=0;i<step.length;i++){
                        console.info(step[i])
                        var tmphtml = '<li>' + step[i] + '</li>'
                        thishtml = thishtml + tmphtml
                    }
                    thishtml = thishtml + '</ol>' + "\n" +
                        ' <hr>' + "\n" +
                        ' <strong><i class="fa fa-pencil margin-r-5"></i> Sql Plan</strong>' ;
                    var isSql = data.SqlDbname.length ;
                    if (isSql != 0) {
                        console.info("no db sql")
                        //获取sql计划
                        thishtml = thishtml + '<ol>' ;
                        for (var i=0;i<data.SqlDbname.length;i++){
                            var sqlMap = data.SqlDbname[i]
                            thishtml = thishtml + '<li class="deploymen-plan-dbname">' ;
                            for (var dbname in sqlMap){
                                thishtml = thishtml + dbname + '<ul type="">' ;
                                for (var j=0;j< sqlMap[dbname].length;j++) {
                                    thishtml = thishtml  + "\n" +
                                        '<li><span class="label label-warning ">SQL:</span>&nbsp;&nbsp;' + sqlMap[dbname][j] + '</li>' ;
                                }
                                thishtml = thishtml + '</ul>' ;
                            }
                            thishtml = thishtml + '</li>' ;
                        }
                        thishtml = thishtml + '</ol>' + "\n" +
                            ' <hr>' + "\n" +
                            '<strong><i class="fa fa-file-text-o margin-r-5"></i> Notes</strong>' + "\n" +
                            '<p>There is no online editing function. Please check the online information .</p>' + "\n" +
                            '</div>' ;

                    } else {
                        thishtml = thishtml + '<strong><i class="fa fa-pencil margin-r-5"></i> Sql Plan</strong>' + "\n" +
                            '<p>No Sql Plan</p>' + "\n" +
                            ' <hr>' + "\n" +
                            '<strong><i class="fa fa-file-text-o margin-r-5"></i> Notes</strong>' + "\n" +
                            '<p>There is no online editing function. Please check the online information .</p>' + "\n" +
                            '</div>' ;
                    }




                    $("#deployment-job-information").html(thishtml)

                    var jobdtime = jobItemResult.DateTime.split(" ")

                    var planhtml = '<ul class="timeline timeline-inverse">' + "\n" +
                        '<li class="time-label">' + "\n" +
                        '<span class="bg-red">' + jobdtime[0] + '</span>' + "\n" +
                        '</li>' ;

                    planhtml = planhtml + "\n" +
                        '<li>' + "\n" +
                        ' <i class="fa fa-calendar-plus-o bg-yellow-active"></i>'  + "\n" +
                        '<div class="timeline-item">' + "\n" +
                        ' <span class="time"><i class="fa fa-clock-o fa-circle-success "></i>' +  jobdtime[1] + '</span>'  + "\n" +
                        '<h3 class="timeline-header"><a href="#">' + jobItemResult.Title + '</a> </h3>'  + "\n" +
                        ' </div>' + "\n" +
                        '</li>'

                    for (var i=0;i<step.length;i++) {
                        planhtml = planhtml + "\n" +
                            '<li>' + "\n" +
                            ' <i class="fa  fa-play bg-blue"></i>'  + "\n" +
                            '<div class="timeline-item">' + "\n" +
                            ' <span class="time"><i class="fa fa-clock-o fa-circle-success "></i> 12:05</span>'  + "\n" +
                            '<h3 class="timeline-header"><a href="#"><span class="label label-success">' + step[i] + '</span></a> </h3>'  + "\n" +
                            '<div class="timeline-header">' + "\n" +
                            '<a class="btn btn-danger btn-xs">执行</a>' + "\n" +
                            '</div>'  + "\n" +
                            '</div>' + "\n" +
                            ' </li>'
                    }

                    if (isSql != 0) {
                        for (var i=0;i<data.SqlDbname.length;i++){
                            var sqlMap = data.SqlDbname[i]
                            for (var dbname in sqlMap){
                                planhtml = planhtml + '<li>' + "\n" +
                                    ' <i class="fa fa-play bg-red-active"></i>' + "\n" +
                                    '<div class="timeline-item">' + "\n" +
                                    '<span class="time"><i class="fa fa-warning fa-circle-success fa-circle-no"></i> 尚未执行</span>' + "\n" +
                                    '<h3 class="timeline-header"><a href="#"><span class="label label-success">' + dbname + '</span></a> </h3>' + "\n" +
                                    '<div class="timeline-header">' + '<ol >' ;

                                for (var j=0;j< sqlMap[dbname].length;j++) {
                                    planhtml = planhtml  + "\n" +
                                        '<li>' + sqlMap[dbname][j] + '</li>' ;
                                }
                                planhtml = planhtml + '</ol>' + '</div>'  + "\n" +
                                    '<div class="timeline-header">'  + "\n" +
                                    '<a class="btn btn-danger btn-xs">执行</a>' + "\n" +
                                    '</div>' + '</div>' + ' </li>'
                            }
                        }
                    }

                    planhtml = planhtml + '<li class="time-label">' + "\n" +
                        ' <span class="bg-green">'  + "\n" +
                        ' 3 Jan. 2014' + "\n" +
                        ' </span>' + "\n" +
                        '</li>' + "\n" +
                        '<li>' + "\n" +
                        ' <i class="fa fa-stop bg-purple"></i>' + "\n" +
                        ' <div class="timeline-item">' + "\n" +
                        '<span class="time"><i class="fa fa-warning fa-circle-success fa-circle-no"></i> 尚未执行</span>'  + "\n" +
                        '<h3 class="timeline-header"><a href="#"><span class="label label-success">DONE</span></a> </h3>' + "\n" +
                        ' </div>' + "\n" +
                        '</li>' + "\n" +
                        '<li>' + "\n" +
                        ' <i class="fa fa-clock-o bg-gray"></i>'  + "\n" +
                        ' </li>'   + "\n" +
                        '</ul>'


                    $("#deployment-plan").html(planhtml)

                }
            }
        )
    })


});

//函数定义
//待数据加载完后，进行页面初始化
function init() {
    var lastjobtime = $("#list-online-job li:last .job-create-dtime  " ).text();
    $.ajax({
        async: true,
        type:'get',
        url:'/api/deploy/getitemplan?lastjobtime=' + lastjobtime ,
        dataType:'json',
        success:function (data){
            console.info("debug 01");

            var jobItemResult = data;
            var thishtml = ' <div class="box-header with-border">' + "\n" +
                '<h3 class="box-title">' + jobItemResult.Title + ' </h3>' + "\n" +
                '</div>' + "\n" +
                '<div class="box-body" >' + "\n" +
                '<strong><i class="fa fa-book margin-r-5"></i> 上线内容</strong>'   + "\n" +
                '<p class="text-muted">'  + "\n" +
                data.Content + "\n" +
                '</p>' + "\n" +
                '<hr>' + "\n" +
                '<strong><i class="fa fa-sort-amount-asc margin-r-5"></i> 上线步骤</strong>' + "\n" +
                '<ol>'  + "\n" ;
            var step = data.DeploymentStep
            for (var i=0;i<step.length;i++){
                console.info(step[i])
                var tmphtml = '<li>' + step[i] + '</li>'
                thishtml = thishtml + tmphtml
            }
            thishtml = thishtml + '</ol>' + "\n" +
                ' <hr>' + "\n" +
                ' <strong><i class="fa fa-pencil margin-r-5"></i> Sql Plan</strong>' ;
            var isSql = data.SqlDbname.length ;
            if (isSql != 0) {
                console.info("no db sql")
                //获取sql计划
                thishtml = thishtml + '<ol>' ;
                for (var i=0;i<data.SqlDbname.length;i++){
                    var sqlMap = data.SqlDbname[i]
                    thishtml = thishtml + '<li class="deploymen-plan-dbname">' ;
                    for (var dbname in sqlMap){
                        thishtml = thishtml + dbname + '<ul type="">' ;
                        for (var j=0;j< sqlMap[dbname].length;j++) {
                                thishtml = thishtml  + "\n" +
                                    '<li><span class="label label-warning ">SQL:</span>&nbsp;&nbsp;' + sqlMap[dbname][j] + '</li>' ;
                        }
                        thishtml = thishtml + '</ul>' ;
                    }
                    thishtml = thishtml + '</li>' ;
                }
                thishtml = thishtml + '</ol>' + "\n" +
                    ' <hr>' + "\n" +
                    '<strong><i class="fa fa-file-text-o margin-r-5"></i> Notes</strong>' + "\n" +
                    '<p>There is no online editing function. Please check the online information .</p>' + "\n" +
                    '</div>' ;

            } else {
                thishtml = thishtml + '<strong><i class="fa fa-pencil margin-r-5"></i> Sql Plan</strong>' + "\n" +
                    '<p>No Sql Plan</p>' + "\n" +
                    ' <hr>' + "\n" +
                    '<strong><i class="fa fa-file-text-o margin-r-5"></i> Notes</strong>' + "\n" +
                    '<p>There is no online editing function. Please check the online information .</p>' + "\n" +
                    '</div>' ;
            }




            $("#deployment-job-information").html(thishtml)

            var jobdtime = jobItemResult.DateTime.split(" ")

            var planhtml = '<ul class="timeline timeline-inverse">' + "\n" +
                '<li class="time-label">' + "\n" +
                '<span class="bg-red">' + jobdtime[0] + '</span>' + "\n" +
                '</li>' ;

            planhtml = planhtml + "\n" +
                '<li>' + "\n" +
                ' <i class="fa fa-calendar-plus-o bg-yellow-active"></i>'  + "\n" +
                '<div class="timeline-item">' + "\n" +
                ' <span class="time"><i class="fa fa-clock-o fa-circle-success "></i>' +  jobdtime[1] + '</span>'  + "\n" +
                '<h3 class="timeline-header"><a href="#">' + jobItemResult.Title + '</a> </h3>'  + "\n" +
                ' </div>' + "\n" +
                '</li>'

            for (var i=0;i<step.length;i++) {
                planhtml = planhtml + "\n" +
                    '<li>' + "\n" +
                    ' <i class="fa  fa-play bg-blue"></i>'  + "\n" +
                    '<div class="timeline-item">' + "\n" +
                    ' <span class="time"><i class="fa fa-clock-o fa-circle-success "></i> 12:05</span>'  + "\n" +
                    '<h3 class="timeline-header"><a href="#"><span class="label label-success">' + step[i] + '</span></a> </h3>'  + "\n" +
                    '<div class="timeline-header">' + "\n" +
                    '<a class="btn btn-danger btn-xs">执行</a>' + "\n" +
                    '</div>'  + "\n" +
                    '</div>' + "\n" +
                    ' </li>'
            }

            if (isSql != 0) {
                for (var i=0;i<data.SqlDbname.length;i++){
                    var sqlMap = data.SqlDbname[i]
                    for (var dbname in sqlMap){
                        planhtml = planhtml + '<li>' + "\n" +
                            ' <i class="fa fa-play bg-red-active"></i>' + "\n" +
                            '<div class="timeline-item">' + "\n" +
                            '<span class="time"><i class="fa fa-warning fa-circle-success fa-circle-no"></i> 尚未执行</span>' + "\n" +
                            '<h3 class="timeline-header"><a href="#"><span class="label label-success">' + dbname + '</span></a> </h3>' + "\n" +
                            '<div class="timeline-header">' + '<ol >' ;

                        for (var j=0;j< sqlMap[dbname].length;j++) {
                            planhtml = planhtml  + "\n" +
                                '<li>' + sqlMap[dbname][j] + '</li>' ;
                        }
                        planhtml = planhtml + '</ol>' + '</div>'  + "\n" +
                            '<div class="timeline-header">'  + "\n" +
                            '<a class="btn btn-danger btn-xs">执行</a>' + "\n" +
                            '</div>' + '</div>' + ' </li>'
                    }
                }
            }

            planhtml = planhtml + '<li class="time-label">' + "\n" +
                ' <span class="bg-green">'  + "\n" +
                ' 3 Jan. 2014' + "\n" +
                ' </span>' + "\n" +
                '</li>' + "\n" +
                '<li>' + "\n" +
                ' <i class="fa fa-stop bg-purple"></i>' + "\n" +
                ' <div class="timeline-item">' + "\n" +
                '<span class="time"><i class="fa fa-warning fa-circle-success fa-circle-no"></i> 尚未执行</span>'  + "\n" +
                '<h3 class="timeline-header"><a href="#"><span class="label label-success">DONE</span></a> </h3>' + "\n" +
                ' </div>' + "\n" +
                '</li>' + "\n" +
                '<li>' + "\n" +
                ' <i class="fa fa-clock-o bg-gray"></i>'  + "\n" +
                ' </li>'   + "\n" +
                '</ul>'


            $("#deployment-plan").html(planhtml)

            }
        }
        )
    }



$(window).on('load',function(){
    init()
});








