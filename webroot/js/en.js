
var api = new API()

function start() {
    console.log('app start')
    $('#btn_oauth_from').on('click', function() {
        oauth('from')
    });

    $('#btn_oauth_to').on('click', function() {
        oauth('to')
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

function oauth(side) {
    api.oauth(side, {
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
}

$(document).ready(start)
