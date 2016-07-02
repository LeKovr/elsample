/*
    jQuery SAT Plugin
    Copyright (c) 2016 elfire.ru
    Licensed under the MIT license
    Version: 1.4
*/

(function($, undefined) {

  var initializeClock = function (duration, display){
    var timeinterval = setInterval(function(){
      display.text(duration--);
      if(duration < 0){
        clearInterval(timeinterval);
        location.reload(true);
      }
    }, 1000);
  };

  var api = function(methodName, data, formId, cbOK, cbErr, options){
    "use strict";
    options = $.extend(true, {
       statusBlock: '#status',
       errorBlock:  '#errors',
       logBlock:    '#log',
       endPoint:    '/api',
       nameSpace:   'App'
    }, options);
    setForm(formId, false);
    $.jsonRPC.withOptions({
      endPoint: options.endPoint,
      namespace: options.nameSpace,
    }, function() {
      $.jsonRPC.request(methodName, {
        params: data,
        success: function(data) {
          setForm(formId, true);
          if (cbOK) cbOK(data.result);
        },
        error: function(data) {
          setForm(formId, true);
          if (cbErr) {
            cbErr(data.error);
          } else {
            $(options.logBlock).html('Error: ' + data.error.message).addClass('alert-danger');
          }
        }
      });
    });
    return false;
  };

  var setForm = function(formId, enabled){
    if (formId) {
      $.each($(formId)[0].elements, function(k, v) {
        if (enabled) {
          if ($(this).hasClass('tmpDisabled')) $(this).removeAttr('disabled').removeClass('tmpDisabled');
        } else {
          if ($(this).attr('disabled') !== 'disabled') $(this).addClass('tmpDisabled').attr('disabled','disabled');
        }
      });
    }
  };
  var setFinal = function(data, msg, origin){
    if (data.Data) window.ELF.Storage.valueSet("satKey", data.Data);
    if (!msg) msg = 'Поздравляем, авторизация пройдена успешно!';
    $('#log').html(msg).removeClass('alert-info').addClass('alert-success');
    $('a#loglink').attr("href","/var/"+data.Phone+"-activate.txt");
    $('#step0').addClass('hide');
    showDiv('#step', "step3", function() {
      $('a#origlink').attr("href", location.href);
      initializeClock(4, $('#time'));
    });
  };

  $.extend({
    sat: {
      version: '1.0',

      init: function(key, host, path) {
        this.api = host + path;
        this.origin = host;
      	// get status
        //var key=$("body").data("satKey"); // key from local storage
        api('Init', {key: key}, null, 
          function(data) {
            if (data.Code == '1') {
              // не истек таймаут отправленного сообщения
              $('#log').html('До повтора отправки СМС осталось <span id="wait"></span> сек').addClass('alert-warning');
              $('#wait').text(data.Data);
              initializeClock(data.Data, $('#wait'));
              showDiv('#step', "step2", function() {
                if (data.Phone) $('#sentPhone').html('+7 '+data.Phone);
              });
            } else if (data.Code == '3') {
              // Доступ разрешен ранее
              setFinal(data, 'Авторизация подтверждена.', host)
            } else {
              showDiv('#step', "step1", function() {
                $("#phone").mask("(999) 999-99-99");
              });

            }
        }, function(data) {
          $('#log').html('Сервис временно недоступен. ('+data.message+')') //TODO - что делать?  
        },
        { endPoint: this.api }
        );
      },

      clearKey: function(link) {
        window.ELF.Storage.valueSet('satKey', '');
        link.innerHTML = '';
        return false;
      },

      sendPhone: function(form) {
        s=$(form).find('input[name="phone"]').val();
        var host = this.origin;
        api('Phone', {key: s}, form, 
          function(data) {
            if (data.Code == '1') {
              // не истек таймаут отправленного сообщения
              $('#log').html('До повтора отправки СМС осталось <span id="wait"></span> сек').addClass('alert-warning');
              $('#wait').text(data.Data);
              initializeClock(data.Data, $('#wait'));
            }
            showDiv('#step', "step2", function() {
              if (data.Phone) $('#sentPhone').html('+7 '+data.Phone);
            });
        }, function(data) {
          $('#log').html('Сервис временно недоступен. ('+data.message+')') //TODO - что делать?  
        },
        { endPoint: this.api }
        );
        return false;
      },

      sendCode: function(form) {
        s=$(form).find('input[name="code"]').val();
        var host = this.origin;
        $('#log').html('Проверка кода...').removeClass('alert-danger').addClass('alert-info');
        api('Code', {key: s}, form
        , function(data, msg) {
          setFinal(data, msg, host)
        }, function(data) {
          if (data.code == -32020) {
            $('#log').html('Неправильный код').addClass('alert-danger');
          } else {
            $('#log').html('Сервис временно недоступен. ('+data.message+')') //TODO - что делать?  
          }
        },
        { endPoint: this.api }
        );
        return false;
      }
    }
  });

})(jQuery);
