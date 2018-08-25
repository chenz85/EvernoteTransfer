
var api = new API()

function start() {
    console.log('app start')
    $('#btn_oauth').on('click', function() {
        api.oauth({
            success: function(data) {
                console.log('oauth:', data)
                console.log('url:', data.authorization_url)
                if (data.authorization_url) {
                    window.location = data.authorization_url
                }
            },
            fail: function() {
                alert('auth failed')
            }
        })
    });

    $('#btn_user').on('click', function() {
        api.user({
            success: function(data) {
                console.log('user data:', data)
            },
            fail: function() {
                alert('get user data failed')
            }
        })
    })
}

$(document).ready(start)
