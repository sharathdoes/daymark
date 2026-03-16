'use client'

import { useEffect } from 'react'
import { useRouter, useSearchParams } from 'next/navigation'
import { useAuth } from '@/lib/store'
import { getMe, setToken as storeToken } from '@/lib/api'

export default function OAuthCallbackPage() {
  const router = useRouter()
  const searchParams = useSearchParams()
  const { setAuth } = useAuth()

  useEffect(() => {
    const handleCallback = async () => {
      const token = searchParams.get('token')

      if (!token) {
        router.push('/login?error=no_token')
        return
      }

      try {
        // Persist token so getMe (and future requests) can use it
        storeToken(token)
        const user = await getMe(token)
        setAuth(user, token)
        router.push('/')
      } catch (error) {
        console.error('OAuth callback error:', error)
        router.push('/login?error=auth_failed')
      }
    }

    handleCallback()
  }, [searchParams, setAuth, router])

  return (
    <div className="min-h-screen flex items-center justify-center bg-black text-white container-narrow">
      <div className="text-center">
        <h1 className="text-2xl font-bold mb-4">Signing you in...</h1>
        <p className="text-[#666666]">Please wait while we complete your authentication.</p>
      </div>
    </div>
  )
}
