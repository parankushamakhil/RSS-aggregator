import { useState, useEffect } from 'react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Label } from '@/components/ui/label'
import { useAuth } from '../App'
import { LogOut, Plus, Rss } from 'lucide-react'

const Dashboard = () => {
  const [feeds, setFeeds] = useState([])
  const [posts, setPosts] = useState([])
  const [newFeedName, setNewFeedName] = useState('')
  const [newFeedUrl, setNewFeedUrl] = useState('')
  const [loading, setLoading] = useState(false)
  const { logout } = useAuth()

  useEffect(() => {
    fetchFeeds()
    fetchPosts()
  }, [])

  const fetchFeeds = async () => {
    try {
      const response = await fetch('http://localhost:8000/feeds', {
        credentials: 'include'
      })
      if (response.ok) {
        const data = await response.json()
        setFeeds(data || [])
      }
    } catch (error) {
      console.error('Error fetching feeds:', error)
    }
  }

  const fetchPosts = async () => {
    try {
      const response = await fetch('http://localhost:8000/posts', {
        credentials: 'include'
      })
      if (response.ok) {
        const data = await response.json()
        setPosts(data || [])
      }
    } catch (error) {
      console.error('Error fetching posts:', error)
    }
  }

  const addFeed = async (e) => {
    e.preventDefault()
    setLoading(true)

    try {
      const response = await fetch('http://localhost:8000/feeds', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        credentials: 'include',
        body: JSON.stringify({ name: newFeedName, url: newFeedUrl }),
      })

      if (response.ok) {
        setNewFeedName('')
        setNewFeedUrl('')
        fetchFeeds()
        fetchPosts()
      }
    } catch (error) {
      console.error('Error adding feed:', error)
    } finally {
      setLoading(false)
    }
  }

  const deleteFeed = async (feedId) => {
    try {
      const response = await fetch(`http://localhost:8000/feeds/${feedId}`, {
        method: 'DELETE',
        credentials: 'include'
      })

      if (response.ok) {
        fetchFeeds()
        fetchPosts()
      }
    } catch (error) {
      console.error('Error deleting feed:', error)
    }
  }

  return (
    <div className="min-h-screen bg-background">
      <header className="border-b">
        <div className="container mx-auto px-4 py-4 flex justify-between items-center">
          <div className="flex items-center space-x-2">
            <Rss className="h-6 w-6" />
            <h1 className="text-2xl font-bold">RSS Aggregator</h1>
          </div>
          <Button onClick={logout} variant="outline" size="sm">
            <LogOut className="h-4 w-4 mr-2" />
            Logout
          </Button>
        </div>
      </header>

      <main className="container mx-auto px-4 py-8">
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
          {/* Add Feed Form */}
          <div className="lg:col-span-1">
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center space-x-2">
                  <Plus className="h-5 w-5" />
                  <span>Add RSS Feed</span>
                </CardTitle>
                <CardDescription>
                  Add a new RSS feed to your collection
                </CardDescription>
              </CardHeader>
              <CardContent>
                <form onSubmit={addFeed} className="space-y-4">
                  <div className="space-y-2">
                    <Label htmlFor="feedName">Feed Name</Label>
                    <Input
                      id="feedName"
                      type="text"
                      value={newFeedName}
                      onChange={(e) => setNewFeedName(e.target.value)}
                      placeholder="e.g., Tech News"
                      required
                    />
                  </div>
                  <div className="space-y-2">
                    <Label htmlFor="feedUrl">Feed URL</Label>
                    <Input
                      id="feedUrl"
                      type="url"
                      value={newFeedUrl}
                      onChange={(e) => setNewFeedUrl(e.target.value)}
                      placeholder="https://example.com/rss"
                      required
                    />
                  </div>
                  <Button type="submit" className="w-full" disabled={loading}>
                    {loading ? 'Adding...' : 'Add Feed'}
                  </Button>
                </form>
              </CardContent>
            </Card>

            {/* Feeds List */}
            <Card className="mt-6">
              <CardHeader>
                <CardTitle>Your Feeds</CardTitle>
                <CardDescription>
                  {feeds.length} feed{feeds.length !== 1 ? 's' : ''} subscribed
                </CardDescription>
              </CardHeader>
              <CardContent>
                <div className="space-y-2">
                  {feeds.length === 0 ? (
                    <p className="text-muted-foreground text-sm">No feeds added yet</p>
                  ) : (
                    feeds.map((feed) => (
                      <div key={feed.id} className="flex items-center justify-between p-2 border rounded">
                        <div>
                          <p className="font-medium">{feed.name}</p>
                          <p className="text-sm text-muted-foreground">{feed.url}</p>
                        </div>
                        <Button
                          onClick={() => deleteFeed(feed.id)}
                          variant="destructive"
                          size="sm"
                        >
                          Delete
                        </Button>
                      </div>
                    ))
                  )}
                </div>
              </CardContent>
            </Card>
          </div>

          {/* Posts List */}
          <div className="lg:col-span-2">
            <Card>
              <CardHeader>
                <CardTitle>Latest Posts</CardTitle>
                <CardDescription>
                  Recent posts from your RSS feeds
                </CardDescription>
              </CardHeader>
              <CardContent>
                <div className="space-y-4">
                  {posts.length === 0 ? (
                    <p className="text-muted-foreground">No posts available. Add some feeds to get started!</p>
                  ) : (
                    posts.map((post) => (
                      <div key={post.id} className="border-b pb-4 last:border-b-0">
                        <h3 className="font-semibold mb-2">
                          <a 
                            href={post.url} 
                            target="_blank" 
                            rel="noopener noreferrer"
                            className="hover:text-primary"
                          >
                            {post.title}
                          </a>
                        </h3>
                        <p className="text-sm text-muted-foreground">
                          {post.published_at && new Date(post.published_at).toLocaleDateString()}
                        </p>
                      </div>
                    ))
                  )}
                </div>
              </CardContent>
            </Card>
          </div>
        </div>
      </main>
    </div>
  )
}

export default Dashboard

