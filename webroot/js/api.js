function request_callback_fail_failback(err) {
    if (err.errMsg) {
        alert('error: [#{0}] {1}'.format(err.errCode, err.errMsg));
    }
 }

function request(url, data, method, success, fail) {

    if (data && method == 'GET') {
        data = undefined;
    }

    fail = fail || request_callback_fail_failback;

    var callback_success = function (response) { //{data, statusCode, header}
        // console.log(response.header);

        var data = response;
        if (typeof data === 'string') {
            try {
                data = JSON.parse(response)
            } catch (e) {
                var error = {
                    errCode: -1,
                    errMsg: 'invalid response data:' + data,
                };
        
                fail.call(null, error);
                return      
            }
        }

        console.log('[http] response data:', data);
        if (data.errCode === undefined || data.errCode == 0) {
            // 处理成功返回数据
            success && success.call(null, data);
        } else {
            // 处理失败信息
            fail.call(null, data); // { errCode, errMsg }
        }
    }

    var callback_fail = function() {
        var error = {
            errCode: -1,
            errMsg: "IOERROR",
        };

        fail.call(null, error);
    }

    var callback_complete = function() {
        // code here if necessary
    }

    // ref: https://api.jquery.com/jquery.post/
    // var req = $.post(url, JSON.stringify(data), null, 'json');
    var req = $.ajax({
        url: url,
        data: JSON.stringify(data),
        type: 'POST',
        dataType: 'json',
        xhrFields: {
            withCredentials: true,
        },
        success: callback_success,
        fail: callback_fail,
        complete: callback_complete,
    });
}

var SERVER_HOST = '/api/';

/** 向api发送请求 */
function api_request(api, data, req_callback) {
  var url = SERVER_HOST + api;
  req_callback = req_callback || {};
  request(url, data, 'POST', req_callback.success, req_callback.fail);
}

var API = (function() {

    ////////////http request//////////////////////////////

    ///////////////////////////

    function API() {

    }

    // 开始oauth验证
    API.prototype.oauth = function(side, callback) {
        var _this = this;

        api_request('en/oauth', {
            side: side,
        }, {
            success: callback && callback.success,
            fail: callback && callback.fail,
        });
    }

    // 获取用户信息
    API.prototype.user = function(callback) {
        var _this = this;

        api_request('en/user', null, {
            success: callback && callback.success,
            fail: callback && callback.fail,
        });
    }

    return API;
})();




