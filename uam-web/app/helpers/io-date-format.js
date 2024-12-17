import Ember from 'ember';

const formatType = ["yyyy-mm-dd", "yyyy-mm-dd HH:MM:ss", "yyyy年mm月dd日"];

export default Ember.Helper.extend({
    compute: function(param, hash) {
        var date = this.getDate(param);
        var year = date.getFullYear(),
            month = date.getMonth() + 1,
            day = date.getDate(),
            hour = "0" + date.getHours(),
            min = "0" + date.getMinutes(),
            sec = "0" + date.getSeconds();

        var index = formatType.indexOf(hash.formatType);

        if (index === -1) {
            return year + "-" + month + "-" + day;
        }

        if (index === 0) {
            return year + "-" + month + "-" + day;
        }

        if (index === 1) {
            return year + "-" + month + "-" + day + " " + hour.substr(hour.length - 2) + ":" + min.substr(min.length - 2) + ":" + sec.substr(sec.length - 2);
        }

        if (index === 2) {
            return year + '年' + month + '月' + day + '日';
        }
        return "";
    },
    getDate: function(date) {
        if (Ember.isBlank(date)) {
            return new Date();
        } else {
            return new Date(date[0]);
        }
    }
});
