'use client'

import { useEffect, useState } from 'react'
import { useRouter } from 'next/navigation'
import { useAuth } from '@/lib/store'
import { getMe, getQuizHistory } from '@/lib/api'
import { UserQuizResult } from '@/lib/types'
import { Card, CardHeader, CardTitle, CardDescription, CardContent } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { LogOut, User as UserIcon, Mail, Calendar } from 'lucide-react'
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'
import { LinkIcon, Check } from 'lucide-react'

export default function ProfilePage() {
  const router = useRouter()
  const { user, setUser, clearAuth, isAuthenticated } = useAuth()
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState('')
  const [quizHistory, setQuizHistory] = useState<UserQuizResult[]>([])
  const [copiedId, setCopiedId] = useState<number | null>(null)

  useEffect(() => {
    // Fetch user and quiz history
    const loadData = async () => {
      try {
        if (!user) {
          const userData = await getMe()
          setUser(userData)
        }
        
        if (isAuthenticated) {
          const history = await getQuizHistory()
          setQuizHistory(history)
        }
      } catch (err) {
        console.error('Failed to load profile data:', err)
        setError('Failed to load profile.')
      } finally {
        setIsLoading(false)
      }
    }

    loadData()
  }, [user, setUser, isAuthenticated])

  const handleLogout = () => {
    clearAuth()
    router.push('/login')
  }

  return (
    <div className="min-h-screen bg-background flex flex-col">
      <main className="flex-1 p-4 md:p-8">
        <div className="w-full max-w-4xl mx-auto space-y-6">
          <Card>
            <CardHeader className="text-center space-y-4">
              <div className="flex justify-center">
                <Avatar className="h-24 w-24">
                  <AvatarImage
                    src={
                      user?.avatar_url ||
                      `https://api.dicebear.com/7.x/identicon/svg?seed=${encodeURIComponent(
                        user?.email || 'user',
                      )}`
                    }
                    alt={user?.name || user?.email || 'Profile picture'}
                  />
                  <AvatarFallback className="text-2xl">
                    {(user?.name || user?.email || '?').charAt(0).toUpperCase()}
                  </AvatarFallback>
                </Avatar>
              </div>
              <div>
                <CardTitle className="text-2xl font-bold">
                  {isLoading ? 'Loading...' : user?.name || 'User Profile'}
                </CardTitle>
                <CardDescription>
                  {user?.provider === 'email' ? 'Email Account' : `${user?.provider} Account`}
                </CardDescription>
              </div>
            </CardHeader>
            <CardContent className="space-y-6">
              {error && (
                <div className="rounded-md bg-destructive/10 p-3 text-sm text-destructive">
                  {error}
                </div>
              )}

              {!isLoading && user && (
                <div className="space-y-4">
                  <div className="flex items-center space-x-3 text-sm">
                    <Mail className="h-5 w-5 text-muted-foreground" />
                    <span>{user.email || 'No email provided'}</span>
                  </div>
                  
                  {user.name && (
                    <div className="flex items-center space-x-3 text-sm">
                      <UserIcon className="h-5 w-5 text-muted-foreground" />
                      <span>{user.name}</span>
                    </div>
                  )}

                  <div className="flex items-center space-x-3 text-sm">
                    <Calendar className="h-5 w-5 text-muted-foreground" />
                    <span>
                      {user.created_at ? `Joined ${new Date(user.created_at).toLocaleDateString()}` : 'Date joined unknown'}
                    </span>
                  </div>
                </div>
              )}

              <Button 
                variant="destructive" 
                className="w-full md:w-auto mt-6"
                onClick={handleLogout}
              >
                <LogOut className="h-4 w-4 mr-2" />
                Sign Out
              </Button>
            </CardContent>
          </Card>

          {/* Quiz History Section */}
          <Card>
            <CardHeader>
              <CardTitle>My Quizzes</CardTitle>
              <CardDescription>Your past quiz results and scores.</CardDescription>
            </CardHeader>
            <CardContent>
              {isLoading ? (
                <p className="text-sm text-muted-foreground">Loading history...</p>
              ) : quizHistory.length === 0 ? (
                <p className="text-sm text-muted-foreground">You haven&apos;t taken any quizzes yet.</p>
              ) : (
                <div className="border rounded-md">
                  <Table>
                    <TableHeader>
                      <TableRow>
                        <TableHead>Date</TableHead>
                        <TableHead>Score</TableHead>
                        <TableHead>Difficulty</TableHead>
                        <TableHead>Topics</TableHead>
                        <TableHead className="text-right">Action</TableHead>
                      </TableRow>
                    </TableHeader>
                    <TableBody>
                      {quizHistory.map(quiz => (
                        <TableRow key={quiz.id}>
                          <TableCell className="font-medium">
                            {new Date(quiz.created_at).toLocaleDateString()}
                          </TableCell>
                          <TableCell>
                            {quiz.score} / {quiz.total_questions} (
                            {Math.round((quiz.score / quiz.total_questions) * 100)}%)
                          </TableCell>
                          <TableCell className="capitalize">{quiz.difficulty}</TableCell>
                          <TableCell className="truncate max-w-[200px]">
                            {/* Assuming categories is a generic string of topics or IDs */}
                            {quiz.categories || 'General'}
                          </TableCell>
                          <TableCell className="text-right">
                            <Button
                              variant="ghost"
                              size="sm"
                              onClick={() => {
                                const url = `${window.location.origin}/quiz/${quiz.quiz_id}`
                                navigator.clipboard.writeText(url)
                                setCopiedId(quiz.id)
                                setTimeout(() => setCopiedId(null), 2000)
                              }}
                            >
                              {copiedId === quiz.id ? (
                                <Check className="w-4 h-4 text-emerald-500" />
                              ) : (
                                <LinkIcon className="w-4 h-4" />
                              )}
                            </Button>
                          </TableCell>
                        </TableRow>
                      ))}
                    </TableBody>
                  </Table>
                </div>
              )}
            </CardContent>
          </Card>
        </div>
      </main>
    </div>
  )
}
