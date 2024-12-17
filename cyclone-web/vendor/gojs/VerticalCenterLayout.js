/**
 * 自定义的一个垂直居中布局样式
 *
 * Created by zhangrong on 16/6/21.
 */

/**
 * @constructor
 * @extends Layout
 * @class
 * This layout assumes the graph is a chain of Nodes,
 * positioning nodes in horizontal rows back and forth, alternating between left-to-right
 * and right-to-left within the {@link #wrap} limit.
 * {@link #spacing} controls the distance between nodes.
 * <p/>
 * When this layout is the Diagram.layout, it is automatically invalidated when the viewport changes size.
 */
function VerticalCenterLayout() {
  go.Layout.call(this);
  // this.isViewportSized = true;
  this._spacing = new go.Size(30, 30);
}
go.Diagram.inherit(VerticalCenterLayout, go.Layout);

/**
 * @ignore
 * Copies properties to a cloned Layout.
 * @this {VerticalCenterLayout}
 * @param {Layout} copy
 * @override
 */
VerticalCenterLayout.prototype.cloneProtected = function (copy) {
  go.Layout.prototype.cloneProtected.call(this, copy);
  copy._spacing = this._spacing;
};

/**
 * This method actually positions all of the Nodes, assuming that the ordering of the nodes
 * is given by a single link from one node to the next.
 * This respects the {@link #spacing} and {@link #wrap} properties to affect the layout.
 * @this {VerticalCenterLayout}
 * @param {Diagram|Group|Iterable} coll the collection of Parts to layout.
 */
VerticalCenterLayout.prototype.doLayout = function (coll) {

  var isDiagram = false;

  if (coll === this.diagram) {
    coll = this.diagram.nodes;
    isDiagram = true;
  } else if (coll === this.group) {
    coll = this.group.memberParts;
    isDiagram = false;
  }

  var it = coll.iterator;
  var all = [];

  while (it.next()) {
    var node = it.value;
    if (!(node instanceof go.Node)) continue;

    if (isDiagram) {
      if (node.isTopLevel) {
        all.push(node);
      }
    } else {
      all.push(node);
    }

  }

  var x = 0;
  var y = 0;

  var maxWidth = Math.max.apply(Math, all.map(function (node) {
    return node.actualBounds.width;
  }));
  var spacing = this.spacing;

  all.forEach(function (n) {
    var dd = n.actualBounds;

    x = ((maxWidth - dd.width) / 2);

    n.move(new go.Point(x, y));

    y += (dd.height + spacing.height);

  });

};

// Public properties

/**
 * Gets or sets the {@link Size} whose width specifies the horizontal space between nodes
 * and whose height specifies the minimum vertical space between nodes.
 * The default value is 30x30.
 * @name VerticalCenterLayout#spacing
 * @function.
 * @return {Size}
 */
Object.defineProperty(VerticalCenterLayout.prototype, "spacing", {
  get: function () {
    return this._spacing;
  },
  set: function (val) {
    if (!(val instanceof go.Size)) {
      throw new Error("new value for VerticalCenterLayout.spacing must be a Size, not: " + val);
    }

    if (!this._spacing.equals(val)) {
      this._spacing = val;
      this.invalidateLayout();
    }
  }
});
