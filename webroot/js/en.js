
var api = new API()

function start() {
    console.log('app start')
    $('#btn_oauth').on('click', function() {
        api.oauth({
            success: function(data) {
                console.log('oauth:', data)
                if (data.authorization_url) {
                    window.location = data.authorization_url
                }
            },
            fail: function() {
                alert('auth failed')
            }
        })
    })
}

$(document).ready(start)
