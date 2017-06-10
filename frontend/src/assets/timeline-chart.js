
var margin = {top: 20, right: 50, bottom: 30, left: 50},
        ohlcHeight = 400,
        width = 960 - margin.left - margin.right,
        height = ohlcHeight /*- margin.top - margin.bottom */;
        

var parseDate = d3.timeParse("%d-%b-%y");

var x = techan.scale.financetime()
        .range([0, width]);

var y = d3.scaleLinear()
        .range([ohlcHeight, 0]);

var candlestick = techan.plot.candlestick()
        .xScale(x)
        .yScale(y);

var yVolume = d3.scaleLinear()
        .range([y(0), y(0.2)]);

var volume = techan.plot.volume()
        .accessor(candlestick.accessor())   // Set the accessor to a ohlc accessor so we get highlighted bars
        .xScale(x)
        .yScale(yVolume);

var xAxis = d3.axisBottom()
        .scale(x);

var yAxis = d3.axisLeft()
        .scale(y);

var volumeAxis = d3.axisRight(yVolume)
        .ticks(3)
        .tickFormat(d3.format(",.3s"));

var volumeAnnotation = techan.plot.axisannotation()
        .axis(volumeAxis)
        .orient("right")
        .width(35);

var svg = null;


$(document).ready(function(){
    svg = d3.select("div.timeline-chart").append("svg")
        .attr("width", width + margin.left + margin.right)
        .attr("height", height + margin.top + margin.bottom)
        .append("g")
        .attr("transform", "translate(" + margin.left + "," + margin.top + ")");

    d3.csv("assets/data.csv", function(error, data) {
        var accessor = candlestick.accessor();
        var volumeAccessor = volume.accessor();

        data = data.slice(0, 200).map(function(d) {
            return {
                date: parseDate(d.Date),
                open: +d.Open,
                high: +d.High,
                low: +d.Low,
                close: +d.Close,
                volume: +d.Volume
            };
        }).sort(function(a, b) {
            return d3.ascending(accessor.d(a), accessor.d(b));
        });

        svg.append("g")
                .attr("class", "candlestick");

        svg.append("g")
                .attr("class", "x axis")
                .attr("transform", "translate(0," + height + ")");

        svg.append("g")
                .attr("class", "y axis")
                .append("text")
                .attr("transform", "rotate(-90)")
                .attr("y", 6)
                .attr("dy", ".71em")
                .style("text-anchor", "end")
                .text("Price");

        svg.append("g")
                .attr("class", "volume");

        svg.append("g")
            .attr("transform", "translate(" + x(1) + ",0)")
            .attr("class", "volume axis")
            .append("text")
                .attr("transform", "translate(0," + height * 0.8 + "), rotate(-90)")
                .style("text-anchor", "end")
                .attr("dy", "-.71em")
                .text("Volume");

        // Data to display initially
        draw(data.slice(0, data.length-20));
        // Only want this button to be active if the data has loaded
        d3.select("button").on("click", function() { draw(data); }).style("display", "inline");
    });
});


function draw(data) {
    x.domain(data.map(candlestick.accessor().d));
    y.domain(techan.scale.plot.ohlc(data, candlestick.accessor()).domain());

    yVolume.domain(techan.scale.plot.volume(data).domain());

    svg.selectAll("g.candlestick").datum(data).call(candlestick);
    svg.selectAll("g.x.axis").call(xAxis);
    svg.selectAll("g.y.axis").call(yAxis);

    svg.select("g.volume").datum(data).call(volume);
    svg.select("g.volume.axis").call(volumeAxis);
}
