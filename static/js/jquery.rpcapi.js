
// Make array of all form input elements
$.fn.serializeObject = function() {
  "use strict";
  var o = {};
  var a = this.serializeArray();
  $.each(a, function() {
    if (o[this.name]) {
      if (!o[this.name].push) {
        o[this.name] = [o[this.name]];
      }
      o[this.name].push(this.value || '');
    } else if (this.name.indexOf('[]') != -1 ) { // && !$.isArray(this.value)
      o[this.name] = [this.value || ''];
    } else {
      o[this.name] = this.value || '';
    }
  });
  return o;
};

// Disable or enable all form input elements
$.fn.inputsEnable = function(enabled){
  var form = this;
  $.each(form[0].elements, function(k, v) {
    if (enabled) {
      if ($(this).hasClass('tmpDisabled')) $(this).removeAttr('disabled').removeClass('tmpDisabled');
    } else {
      if ($(this).attr('disabled') !== 'disabled') $(this).addClass('tmpDisabled').attr('disabled','disabled');
    }
  });
};


(function($, undefined) {

  // Validate a params hash
  var _validateConfigOptions = function(options) {
    if(options === undefined) return;
    if(options.logBlock && typeof(options.logBlock) !== 'object'){
      throw("logBlock must be an object");
    }
  };

  $.extend({
    rpc: {
      setup: function(options) {
        _validateConfigOptions(options);
        options = $.extend(true, {
          logBlock:    '#log',
          endPoint:    '/api',
          nameSpace:   'App'
        }, options);
        this.options = options;
      },

      api: function(options){
        _validateConfigOptions(options);
        if (this.options === undefined) {
          this.setup();
        }
        options = $.extend(true, this.options, options);
        var $form = $(options.Form);
        var p = (options.Params === undefined) ? $form.serializeObject() : options.Params;
        $form.inputsEnable(false);
        $.jsonRPC.withOptions({
          endPoint: options.endPoint,
          namespace: options.nameSpace,
        }, function() {
          $.jsonRPC.request(options.Method, {
            params: p,
            success: function(data) {
              $form.inputsEnable(true);
              if (options.onSuccess) options.onSuccess(data.result);
            },
            error: function(data) {
              $form.inputsEnable(true);
              if (options.onError) {
                options.onError(data.error);
              } else {
                $(options.logBlock).html('Error: ' + data.error.message).addClass('alert-danger');
              }
            }
          });
        });
        return false;
      }
    }
  });

})(jQuery);
