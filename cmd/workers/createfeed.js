const FEEDS = [
  {
    Name: "The Hindu",
    URL: "https://www.thehindu.com/feeder/default.rss",
    categoryIds: [3],
  },
  {
    Name: "The Hindu",
    URL: "https://www.thehindu.com/news/national/feeder/default.rss",
    categoryIds: [23],
  },
  {
    Name: "The Hindu",
    URL: "https://www.thehindu.com/news/international/feeder/default.rss",
    categoryIds: [4],
  },
  {
    Name: "The Hindu",
    URL: "https://www.thehindu.com/business/Economy/feeder/default.rss",
    categoryIds: [10, 11],
  },
  {
    Name: "The Hindu",
    URL: "https://www.thehindu.com/business/Industry/feeder/default.rss",
    categoryIds: [10],
  },
  {
    Name: "The Hindu",
    URL: "https://www.thehindu.com/sport/cricket/feeder/default.rss",
    categoryIds: [1, 5],
  },
  {
    Name: "The Hindu",
    URL: "https://www.thehindu.com/sport/football/feeder/default.rss",
    categoryIds: [1, 21],
  },
  {
    Name: "The Hindu",
    URL: "https://www.thehindu.com/sci-tech/science/feeder/default.rss",
    categoryIds: [13, 2],
  },
  {
    Name: "The Hindu",
    URL: "https://www.thehindu.com/sci-tech/technology/feeder/default.rss",
    categoryIds: [13, 6],
  },
  {
    Name: "The Hindu",
    URL: "https://www.thehindu.com/sci-tech/agriculture/feeder/default.rss",
    categoryIds: [13, 12],
  },
  {
    Name: "The Hindu",
    URL: "https://www.thehindu.com/news/cities/Delhi/feeder/default.rss",
    categoryIds: [16, 14],
  },
  {
    Name: "The Hindu",
    URL: "https://www.thehindu.com/news/cities/Mumbai/feeder/default.rss",
    categoryIds: [16, 15],
  },
  {
    Name: "The Hindu",
    URL: "https://www.thehindu.com/news/cities/Hyderabad/feeder/default.rss",
    categoryIds: [16, 18],
  },
  {
    Name: "The Hindu",
    URL: "https://www.thehindu.com/news/cities/Chennai/feeder/default.rss",
    categoryIds: [16, 19],
  },
  {
    Name: "The Hindu",
    URL: "https://www.thehindu.com/news/cities/bangalore/feeder/default.rss",
    categoryIds: [16, 20],
  },
  {
    Name: "The Hindu",
    URL: "https://www.thehindu.com/news/cities/Kochi/feeder/default.rss",
    categoryIds: [16],
  },
  {
    Name: "The Hindu",
    URL: "https://www.thehindu.com/news/cities/Vijayawada/feeder/default.rss",
    categoryIds: [16],
  },
  {
    Name: "The Hindu",
    URL: "https://www.thehindu.com/news/cities/Visakhapatnam/feeder/default.rss",
    categoryIds: [16],
  },
  {
    Name: "Indian Express",
    URL: "https://indianexpress.com/feed/",
    categoryIds: [23, 3],
  },
  {
    Name: "Indian Express",
    URL: "https://indianexpress.com/section/world/feed/",
    categoryIds: [4],
  },
  {
    Name: "Indian Express",
    URL: "https://indianexpress.com/section/business/economy/feed/",
    categoryIds: [10],
  },
  {
    Name: "Indian Express",
    URL: "https://indianexpress.com/section/business/banking-and-finance/feed/",
    categoryIds: [10, 11],
  },
  {
    Name: "Indian Express",
    URL: "https://indianexpress.com/section/entertainment/feed/",
    categoryIds: [22],
  },
  {
    Name: "Indian Express",
    URL: "https://indianexpress.com/section/sports/feed/",
    categoryIds: [1],
  },
  {
    Name: "Indian Express",
    URL: "https://indianexpress.com/section/technology/feed/",
    categoryIds: [13, 6],
  },
  {
    Name: "Indian Express",
    URL: "https://indianexpress.com/section/cities/delhi/feed/",
    categoryIds: [16, 14],
  },
  {
    Name: "Indian Express",
    URL: "https://indianexpress.com/section/cities/mumbai/feed/",
    categoryIds: [16, 15],
  },
  {
    Name: "Indian Express",
    URL: "https://indianexpress.com/section/cities/bangalore/feed/",
    categoryIds: [16, 20],
  },
  {
    Name: "Indian Express",
    URL: "https://indianexpress.com/section/cities/pune/feed/",
    categoryIds: [16],
  },
  {
    Name: "Indian Express",
    URL: "https://indianexpress.com/section/cities/ahmedabad/feed/",
    categoryIds: [16],
  },
  {
    Name: "Indian Express",
    URL: "https://indianexpress.com/section/cities/lucknow/feed/",
    categoryIds: [16],
  },
  {
    Name: "Indian Express",
    URL: "https://indianexpress.com/section/cities/chandigarh/feed/",
    categoryIds: [16],
  },
  {
    Name: "Indian Express",
    URL: "https://indianexpress.com/section/cities/kolkata/feed/",
    categoryIds: [16],
  },
  {
    Name: "Times of India",
    URL: "https://timesofindia.indiatimes.com/rssfeedstopstories.cms",
    categoryIds: [3],
  },
  {
    Name: "Times of India",
    URL: "https://timesofindia.indiatimes.com/rssfeeds/-2128936835.cms",
    categoryIds: [23],
  },
  {
    Name: "Times of India",
    URL: "https://timesofindia.indiatimes.com/rssfeeds/296589292.cms",
    categoryIds: [4],
  },
  {
    Name: "Times of India",
    URL: "https://timesofindia.indiatimes.com/rssfeeds/303594735.cms",
    categoryIds: [4, 24],
  },
  {
    Name: "Times of India",
    URL: "https://timesofindia.indiatimes.com/rssfeeds/669432.cms",
    categoryIds: [4],
  },
  {
    Name: "Times of India",
    URL: "https://timesofindia.indiatimes.com/rssfeeds/1898055.cms",
    categoryIds: [10],
  },
  {
    Name: "Times of India",
    URL: "https://timesofindia.indiatimes.com/rssfeeds/1898274.cms",
    categoryIds: [10, 25],
  },
  {
    Name: "Times of India",
    URL: "https://timesofindia.indiatimes.com/rssfeeds/4719148.cms",
    categoryIds: [1, 5],
  },
  {
    Name: "Times of India",
    URL: "https://timesofindia.indiatimes.com/rssfeeds/4719161.cms",
    categoryIds: [1, 21],
  },
  {
    Name: "Times of India",
    URL: "https://timesofindia.indiatimes.com/rssfeeds/3942921.cms",
    categoryIds: [13, 6],
  },
  {
    Name: "Times of India",
    URL: "https://timesofindia.indiatimes.com/rssfeeds/26435941.cms",
    categoryIds: [9],
  },
  {
    Name: "Times of India",
    URL: "https://timesofindia.indiatimes.com/rssfeeds/-2128672765.cms",
    categoryIds: [13, 2],
  },
  {
    Name: "Times of India",
    URL: "https://timesofindia.indiatimes.com/rssfeeds/4473.cms",
    categoryIds: [16, 14],
  },
  {
    Name: "Times of India",
    URL: "https://timesofindia.indiatimes.com/rssfeeds/4481.cms",
    categoryIds: [16, 15],
  },
  {
    Name: "Times of India",
    URL: "https://timesofindia.indiatimes.com/rssfeeds/3947.cms",
    categoryIds: [16, 20],
  },
  {
    Name: "Times of India",
    URL: "https://timesofindia.indiatimes.com/rssfeeds/3941.cms",
    categoryIds: [16, 18],
  },
  {
    Name: "Times of India",
    URL: "https://timesofindia.indiatimes.com/rssfeeds/2950627.cms",
    categoryIds: [16, 19],
  },
  {
    Name: "Times of India",
    URL: "https://timesofindia.indiatimes.com/rssfeeds/3950.cms",
    categoryIds: [16],
  },
  {
    Name: "Times of India",
    URL: "https://timesofindia.indiatimes.com/rssfeeds/3951.cms",
    categoryIds: [16],
  },
  {
    Name: "Times of India",
    URL: "https://timesofindia.indiatimes.com/rssfeeds/3942.cms",
    categoryIds: [16],
  },
  {
    Name: "Times of India",
    URL: "https://timesofindia.indiatimes.com/rssfeeds/3943.cms",
    categoryIds: [16],
  },
  {
    Name: "Times of India",
    URL: "https://timesofindia.indiatimes.com/rssfeeds/3952.cms",
    categoryIds: [16],
  },
  {
    Name: "Times of India",
    URL: "https://timesofindia.indiatimes.com/rssfeeds/4118245.cms",
    categoryIds: [16],
  },
  {
    Name: "Times of India",
    URL: "https://timesofindia.indiatimes.com/rssfeeds/4118252.cms",
    categoryIds: [16],
  },
  {
    Name: "Times of India",
    URL: "https://timesofindia.indiatimes.com/rssfeeds/4118266.cms",
    categoryIds: [16],
  },
  {
    Name: "Times of India",
    URL: "https://timesofindia.indiatimes.com/rssfeeds/3948.cms",
    categoryIds: [16],
  },
  {
    Name: "Times of India",
    URL: "https://timesofindia.indiatimes.com/rssfeeds/4118248.cms",
    categoryIds: [16],
  },
  {
    Name: "Times of India",
    URL: "https://timesofindia.indiatimes.com/rssfeeds/4118258.cms",
    categoryIds: [16],
  },
];

async function createFeeds() {
  const url = "http://localhost:8080/feed/create";

  console.log(`🚀 Starting feed creation for ${FEEDS.length} sources...`);

  for (const feed of FEEDS) {
    try {
      const response = await fetch(url, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(feed),
      });

      if (response.ok) {
        console.log(`✅ Success: ${feed.Name} - ${feed.URL.split("/").pop()}`);
      } else {
        const errorText = await response.text();
        console.error(
          `❌ Failed: ${feed.Name}. Status: ${response.status}. Error: ${errorText}`,
        );
      }
    } catch (error) {
      console.error(`💥 Network Error for ${feed.Name}:`, error.message);
    }
  }

  console.log("✨ All requests processed.");
}

// createFeeds();

async function checkFeeds() {
  for (const feed of FEEDS) {
    try {
      const response = await fetch(
        `http://localhost:8080/debug-rss?url=${encodeURIComponent(feed.URL)}`,
      );

      const data = await response.json();

      console.log(`[${data.statusCode}] ${feed.Name} -> ${feed.URL}`);
    } catch (err) {
      console.log(`[ERROR] ${feed.Name} -> ${feed.URL}`, err.message);
    }
  }
}

checkFeeds();