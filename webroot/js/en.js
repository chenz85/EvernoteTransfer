
var api = new API()

function start() {
    console.log('app start')
    $('#btn_oauth_from_en').on('click', function() {
        oauth('from', 'en')
    });
    $('#btn_oauth_from_yx').on('click', function() {
        oauth('from', 'yx')
    });

    $('#btn_oauth_to_en').on('click', function() {
        oauth('to', 'en')
    });
    $('#btn_oauth_to_yx').on('click', function() {
        oauth('to', 'yx')
    });

    $('#btn_transfer').on('click', function() {
        api.xfer({
            success: function(data) {
                console.log('tx id:', data)
            },
            fail: function() {
                alert('transfer failed')
            }
        })
    })
}

function oauth(side, svr) {
    api.oauth(side, svr, {
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
