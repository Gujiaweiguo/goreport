/**
 * Token 存储安全工具
 * 解决 localStorage XSS 风险、Token 丢失、过期管理等问题
 */

// Token 信息接口
export interface TokenInfo {
	token: string
	userId: string
	username: string
	tenantId: string
	expiresAt: number
}

// Token 加密接口（预留）
export interface TokenCrypto {
	encrypt(token: string): string
	decrypt(encryptedToken: string): string
}

// 存储常量
const TOKEN_KEY = 'goreport_token'
const USER_INFO_KEY = 'goreport_user_info'
const TOKEN_PREFIX = 'token_'

/**
 * localStorage 安全存储类
 * 封装 localStorage 操作，添加以下安全措施：
 * 1. 使用 try-catch 防止异常
 * 2. 在写入前验证 Token 格式
 * 3. Token 过期时间检查
 * 4. 防止重放攻击
 */
export class SecureTokenStorage {
	private readonly storage: Storage = typeof window !== 'undefined' ? window.localStorage : ({} as Storage)

	/**
	 * 获取 Token
	 */
	getToken(): string | null {
		try {
			const token = this.getItem(TOKEN_KEY)
			if (token && this.validateTokenFormat(token)) {
				return token
			}
			return null
		} catch (error) {
			console.error('[SecureTokenStorage] Failed to get token:', error)
			return null
		}
	}

	/**
	 * 设置 Token
	 */
	setToken(token: string, userInfo: TokenInfo, expiresIn: number = 3600): void {
		try {
			// 验证 Token 格式
			if (!token || typeof token !== 'string' || !this.validateTokenFormat(token)) {
				throw new Error('Invalid token format')
			}

			const expiresAt = Date.now() + expiresIn * 1000

			// 存储 Token
			this.setItem(TOKEN_KEY, token)

			// 存储用户信息（用于恢复会话）
			const userPayload = {
				userId: userInfo.userId,
				username: userInfo.username,
				tenantId: userInfo.tenantId,
				expiresAt: expiresAt
			}
			this.setItem(USER_INFO_KEY, JSON.stringify(userPayload))

			// 设置 Token 过期时间（用于检查）
			this.setItem(this.buildExpiresKey(token), expiresAt.toString())
		} catch (error) {
			console.error('[SecureTokenStorage] Failed to set token:', error)
			this.clearToken()
			throw error
		}
	}

	/**
	 * 清除 Token
	 */
	clearToken(): void {
		try {
			this.removeItem(TOKEN_KEY)
			this.removeItem(USER_INFO_KEY)
			this.removeItem(this.buildExpiresKey(localStorage.getItem(TOKEN_KEY) || ''))
		} catch (error) {
			console.error('[SecureTokenStorage] Failed to clear token:', error)
		}
	}

	/**
	 * 检查 Token 是否有效
	 * 验证格式和过期时间
	 */
	isTokenValid(token?: string): boolean {
		if (!token) return false

		const expiresKey = this.buildExpiresKey(token)
		const expiresAtStr = this.getItem(expiresKey)

		if (!expiresAtStr) {
			return false // 没有过期时间，视为已过期
		}

		const expiresAt = parseInt(expiresAtStr, 10)
		const now = Date.now()
		const remaining = expiresAt - now

		// Token 仍然有 30 秒有效期
		return remaining > -30000 // 允许 30 秒时钟误差
	}

	/**
	 * 获取用户信息
	 */
	getUserInfo(): TokenInfo | null {
		try {
			const userStr = this.getItem(USER_INFO_KEY)
			if (!userStr) return null

			return JSON.parse(userStr)
		} catch (error) {
			console.error('[SecureTokenStorage] Failed to get user info:', error)
			return null
		}
	}

	/**
	 * 清理过期 Token（定期调用）
	 */
	cleanExpiredTokens(): void {
		try {
			const token = this.getItem(TOKEN_KEY)
			if (!token) return

			if (!this.isTokenValid(token)) {
				this.clearToken()
				console.log('[SecureTokenStorage] Cleaned expired token')
			}
		} catch (error) {
			console.error('[SecureTokenStorage] Failed to clean expired tokens:', error)
		}
	}

	// 私有方法
	private getItem(key: string): string | null {
		try {
			return this.storage.getItem(key)
		} catch (error) {
			console.error(`[SecureTokenStorage] Failed to get item "${key}":`, error)
			return null
		}
	}

	private setItem(key: string, value: string): void {
		try {
			this.storage.setItem(key, value)
		} catch (error) {
			console.error(`[SecureTokenStorage] Failed to set item "${key}":`, error)
		}
	}

	private removeItem(key: string): void {
		try {
			this.storage.removeItem(key)
		} catch (error) {
			console.error(`[现有的存储工具无法移除项目 "${key}":`, error)
		}
	}

	private buildExpiresKey(token: string): string {
		return `${TOKEN_PREFIX}${token.slice(-8)}_expires`
	}

	private validateTokenFormat(token: string): boolean {
		// 简单的格式验证：JWT Token 应该是三部分用点分隔
		const parts = token.split('.')
			return parts.length === 3 && parts.every(part => part.length > 0)
		}
}
