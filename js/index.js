$(function(){ 
	var name = getUrlParam('name');
	if (name != null && name !='') {
		$('#input_loginname').val(name);
	}
	$('#modal_login').modal('show');
}); 

var addr = '127.0.0.1:10328';
var wsaddr = 'ws://'+addr+'/login'
var ws;
var lasttime = '';

$('#button_login').click(function (){
	var name = $('#input_loginname').val();
	if (name == '' ) {
		alert('名字不能为空!');
		return false;
	} else if ( name.length > 10) {
		alert('名字不能超过10个字符');
		return false;
	}

	if (ws) {
		return false;
	}

	ws = new WebSocket(wsaddr+'?name='+name);
	ws.onopen = function(evt) {
		console.log("CONNECT !")
		$('#modal_login').modal('hide');
	}

	ws.onmessage = function(evt) {
		newmsg = dealMsg(evt.data)
		$('.all-msg').append(newmsg);
		$('.all-msg').scrollTop( $('.all-msg')[0].scrollHeight );

	}

	ws.onclose = function(evt) {
		console.log("CONNECT CLOSE")
		var msg = '<div class="row time-style" >连接已断开,请刷新</div >';
		$('.all-msg').append(msg);
		$('.all-msg').scrollTop( $('.all-msg')[0].scrollHeight );
		ws = null;
	}

	ws.onerror = function(evt) {
		console.log("ERROR: " + evt.data);
	}

	getHistory();
});


$('#button_send').click(function (){
	var msg = $('#input_msg').val();
	if (msg == '') {
		alert('消息不能为空!');
		return ;
	}
	if (!ws) {
		return false;
	}
	ws.send(msg);
	$('#input_msg').val('');
});


$('#input_msg').keyup(function (e){
	var ev = window.event||e;
	if (ev.keyCode == 13) {
		$('#button_send').click();
	}

});

function getHistory() {
	$.ajax({   
		url:'http://'+addr+'/history',
		type:'get',
		dataType:'JSON',
		data:'',
		async : true,
		timeout : 300,
		error:function(data){   
			console.log('ajax fail:', data);
		},   
		success:function(data){   
			var oldmsg = '';
			for (i=0;i< data.length;i++) {
				var recvdata = JSON.stringify(data[i]);
				oldmsg += dealMsg(recvdata);
			}
			if (oldmsg != '') {
				oldmsg += '<div class="row time-style" >以上为历史消息</div >';
				$('.all-msg').prepend(oldmsg);
				$('.all-msg').scrollTop( $('.all-msg')[0].scrollHeight );
			}
		}
	});
}

function dealMsg(recvdata) {
	var newmsg = '';
	var obj = JSON.parse(recvdata);

	var msg = obj['msg'];
	var photo = obj['photo'];
	var time = obj['time'];
	var name = obj['name'];
	var alertmsg = name + ' 发表于 ' + time;
	msg = msg.replace(/\n/g,'<br>');


	var needtime = false;
	if (lasttime == '') {
		needtime = true;
	} else if (time.substr(0,15) > lasttime.substr(0,15) ) {
		needtime = true;
	} else if (parseInt(time[15]) >= parseInt(lasttime[15]) + 5) {
		needtime = true;
	}

	if (needtime) {
		newmsg += '<div class="row time-style" > '+ time.substr(11,5) +' </div >';
		lasttime = time;
	}

	newmsg += '<div class="row single-msg" >';
	newmsg += '	<img class="img-style" width="40" height="40" src="'+photo+'" onclick="alert(\''+alertmsg+'\');" />';
	newmsg += '	<div class="msg-style">'+msg+'</div>';
	newmsg += '</div>';
	return newmsg;
}

function addZero(s) {
	return s < 10 ? '0' + s: s;
}

function getCurrentTime() {
	var myDate = new Date();
	var year=myDate.getFullYear();
	var month=myDate.getMonth()+1;
	var date=myDate.getDate(); 
	var h=myDate.getHours();	   //获取当前小时数(0-23)
	var m=myDate.getMinutes();	 //获取当前分钟数(0-59)
	var s=myDate.getSeconds();  

	var now=year+'-'+addZero(month)+"-"+addZero(date)+" "+addZero(h)+':'+addZero(m)+":"+addZero(s);
	return now;
}

function getUrlParam(name) {
	var reg = new RegExp("(^|&)" + name + "=([^&]*)(&|$)"); //构造一个含有目标参数的正则表达式对象
	var r = window.location.search.substr(1).match(reg);  //匹配目标参数
	//if (r != null) return unescape(r[2]); return null; //返回参数值
	if (r != null) return decodeURI(r[2]); return null; //返回参数值
}


/*
window.onfocus = function() {
	console.log(getCurrentTime());
};

*/
