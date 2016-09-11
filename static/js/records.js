      var params;

    function sendLogout() {
      $.cookie('elfire_sso_token', null, { path: '/' });
      location.reload(true);
      return false;
    }

    function getRows(params) {
        var by = params['by'];
        if (!by) by = 10;

      $('#log').html('Загрузка данных...').removeClass('alert-danger').addClass('alert-info');
      $.rpc.api({
        Form:      $('#query'),
        nameSpace: 'Service',
        Method:    'List',
        endPoint:  '/my/api/v1',
        Params:    params,
        onSuccess: function(result) {
          $('#log').html('');
          if (result.length < params['by'] || params['by'] == -1) {
            // last page => hide "Next" button
            $('#next').addClass('hide');
            if (result.length == 0) {
              $('#log').html('Нет данных');
            }
          }
          $.each(result, function(index, item) {
            var p = $('<tr>');
            p.append($('<td>').text(moment(item.Stamp).format("L LT")));
            p.append($('<td>').text(item.Phone));
            p.append($('<td>').text(item.IP));
            p.append($('<td>').text(item.Status));
            $('#data > tbody:last-child').append(p);
          });
          params["offset"] += result.length;
        },
        onError: function(error) {
          $('#log').html('Ошибка: ' + error.message).removeClass('alert-info').addClass('alert-danger');
        }
      });
      return false;
    }

    function sendNext() {
      getRows(params);
    }

    $.ajaxSetup({

      'beforeSend': function(xhr) {
        xhr.setRequestHeader('X-ELFIRE-Token', 'Bearer ' + $.cookie('elfire_sso_token'));
      }
    });

    $(function () {
      //Инициализация
      $("#dtAfter").datetimepicker();
      $("#dtBefore").datetimepicker();

      //При изменении даты в dtAfter, она устанавливается как минимальная для dtBefore
      $("#dtAfter").on("dp.change",function (e) {
        $("#dtBefore").data("DateTimePicker").minDate(e.date);
      });
      //При изменении даты в dtBefore, она устанавливается как максимальная для 8dtAfter
      $("#dtBefore").on("dp.change",function (e) {
        $("#dtAfter").data("DateTimePicker").maxDate(e.date);
      });

      // http://stackoverflow.com/questions/5448545/how-to-retrieve-get-parameters-from-javascript
      params = decodeURIComponent(window.location.search.slice(1))
        .split('&')
        .reduce(function _reduce (/*Object*/ a, /*String*/ b) {
          b = b.split('=');
          if ($.isArray(a[b[0]])) {
            a[b[0]].push(b[1])
          } else if (typeof a[b[0]] === 'undefined' && b[1] != undefined) {
            a[b[0]] = decodeURIComponent(b[1].replace(/\+/g, '%20'));
          } else if (b[1] != undefined) {
            a[b[0]] = [a[b[0]], b[1]]
          }
          return a;
        }, {});
      params["by"] = (params["by"] === undefined) ? 10 : parseInt(params["by"]);
      params["offset"] = 0;

      var $f = $('#query');
      for (var param in params) {
        field = $f.find('input[name='+param+']')[0] || $f.find('textarea[name='+param+']')[0] || $f.find('select[name='+param+']')[0];
          if (field === undefined) {
            continue;
          } else if (field.type === 'select') {
            var fieldvalue = params[param];
            $f.find('select[name='+param+']')[0].val( fieldvalue ).change();
          } else {
            field.value = params[param];
          }
      }
      moment.fn.toJSON = function() { return this.format(); } // keep timezone

      if (params["before"] == undefined || params["before"] == "") {
        delete params["before"];
      } else {
        params["before"] = moment(params["before"],"L LT").toJSON();
      }
      if (params["after"] == undefined || params["after"] == "") {
        delete params["after"];
      } else {
        params["after"] = moment(params["after"],"L LT").toJSON();
      }

      //if (params['phone'] || params['ip'] || params['before'] || params['after']) {
        $('#data').removeClass('hide');
        getRows(params);
        $('#next').removeClass('hide');
      //}
    });

