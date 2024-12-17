import Ember from 'ember';

/**
 * Alert Component
 ```html
 ``` 
 */

export default Ember.Component.extend({
	/**
	 * [tagName description]
	 */
	tagName: 'div',
	attributeBindings: ['disabled', 'role'],
	classNames: 'io-tooltips',
	classNamePrefix: 'io-tooltips-',
	classNameBindings: ['typeClass','displayClass','activeClass'],
	/**
	 * @attribute triggerEvent
	 * @type {String [mosueover | click]}
	 */
	triggerEvent: 'mouseover',
	_hidden: true,
	showTips:true,
	_style:function(){
		var getStyle='';
		if (!this.get('_hidden')) {
			var dw = this.$('.io-tooltips-title').width();
			var dh = this.$('.io-tooltips-title').height();
			var tw = this.$('.io-tooltips-content').width();
			var th = this.$('.io-tooltips-content').height();
			// console.log("dw:"+dw+" dh:"+dh);
			var topY    = -10-th;
			var bottomY = dh+10;
			var leftX   = -5-tw;
			var rightX  = dw+20;

			var topLeft =  (dw/2-tw-10)/2;
			var topCenter = (dw-tw)/2;
			var topRight  = (dw/2+dw-tw+10)/2;

			var bottomLeft  = (dw/2-tw-10)/2;
			var bottomCenter  = (dw-tw)/2;
			var bottomRight  = (dw/2+dw-tw+10)/2;

			var leftTop    = (dh/2-th-10)/2;
			var leftCenter = (dh-th-10)/2;
			var leftBottom = (dh/2+dh-th-10)/2;

			var rightTop     = (dh/2-th-10)/2;
			var rightCenter  = (dh-th-10)/2;
			var rightBottom  = (dh/2+dh-th-10)/2;

			var type=this.get('type');
			
			switch (type) {
				case 'topLeft':
					getStyle="left:"+topLeft+"px;top:"+topY+"px;";break;
				case 'top':
					getStyle="left:"+topCenter+"px;top:"+topY+"px;";break;
				case 'topRight':
					getStyle="left:"+topRight+"px;top:"+topY+"px;";break;
				case 'rightTop':
					getStyle="left:"+rightX+"px;top:"+rightTop+"px;";break;
				case 'right':
					getStyle="left:"+rightX+"px;top:"+rightCenter+"px;";break;
				case 'rightBottom':
					getStyle="left:"+rightX+"px;top:"+rightBottom+"px;";break;
				case 'bottomRight':
					getStyle="left:"+bottomRight+"px;top:"+bottomY+"px;";break;
				case 'bottom':
					getStyle="left:"+bottomCenter+"px;top:"+bottomY+"px;";break;
				case 'bottomLeft':
					getStyle="left:"+bottomLeft+"px;top:"+bottomY+"px;";break;
				case 'leftBottom':
					getStyle="left:"+leftX+"px;top:"+leftBottom+"px;";break;
				case 'left':
					getStyle="left:"+leftX+"px;top:"+leftCenter+"px;";break;
				case 'leftTop':
					getStyle="left:"+leftX+"px;top:"+leftTop+"px;";break;
				default:
					getStyle="left:"+topCenter+"px;top:"+topY+"px;";break;
			};
		}
		return getStyle;
	}.property('_hidden'),
	/**
	 * classNameBindings
	 * @return {[type]} [description]
	 */
	typeClass: function() {
			return this.get('classNamePrefix') + "placement-" +this.get('type');
	}.property('type'),
	displayClass: function() {
		if (this.get('_hidden')) {
			return this.get('classNamePrefix') + 'hidden';
		} 
		return '';
	}.property('_hidden'),
	activeClass: function() {
		if (!this.get('_hidden')) {
			return this.get('classNamePrefix') + 'active';
		} else {
			return '';
		}
	}.property('_hidden'),
	/**
	 * mosueover event description
	 * @type {[type]}
	 */
	t: null,
	_mouseover: false,
	mouseEnter() {
		//console.log("悬浮");
		if (this.get('triggerEvent') !== 'mouseover') {
			return
		}
		clearTimeout(this.get('t'));
		this.set('_mouseover', true);
		this.set('_hidden', false);
		this.set('showTips',true);
	},
	mouseLeave() {
		//console.log("远离");
		if (this.get('triggerEvent') !== 'mouseover') {
			return
		}
		const t = setTimeout(function() {
			this.set('_mouseover', false);
			this.set('_hidden', true);
			this.set('showTips',false);
		}.bind(this), 200);
		this.set('t', t);
	},
});