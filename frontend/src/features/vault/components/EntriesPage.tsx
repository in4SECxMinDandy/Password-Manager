import { useState } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import {
  ArrowLeft,
  Plus,
  Loader2,
  MoreHorizontal,
  Trash2,
  Edit,
  Copy,
  Star,
  Search,
  LogIn,
  CreditCard,
  FileText,
  User,
} from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import { toast } from 'sonner'
import { apiClient } from '@/lib/api-client'
import type { VaultEntry } from '@/lib/api-types'

const entryTypeIcons: Record<string, React.ReactNode> = {
  login: <LogIn className="h-4 w-4" />,
  card: <CreditCard className="h-4 w-4" />,
  note: <FileText className="h-4 w-4" />,
  identity: <User className="h-4 w-4" />,
}

export function EntriesPage() {
  const { vaultId } = useParams<{ vaultId: string }>()
  const navigate = useNavigate()
  const queryClient = useQueryClient()
  const [isCreateOpen, setIsCreateOpen] = useState(false)
  const [newEntryTitle, setNewEntryTitle] = useState('')
  const [newEntryType, setNewEntryType] = useState<VaultEntry['type']>('login')
  const [searchQuery, setSearchQuery] = useState('')

  const { data: vault } = useQuery({
    queryKey: ['vault', vaultId],
    queryFn: () => apiClient.getVault(vaultId!),
    enabled: !!vaultId,
  })

  const { data, isLoading } = useQuery({
    queryKey: ['entries', vaultId],
    queryFn: () => apiClient.getEntries(vaultId!),
    enabled: !!vaultId,
  })

  const createMutation = useMutation({
    mutationFn: (entry: { title: string; type: VaultEntry['type']; data: string }) =>
      apiClient.createEntry(vaultId!, {
        ...entry,
        data: JSON.stringify({}),
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['entries', vaultId] })
      setIsCreateOpen(false)
      setNewEntryTitle('')
      setNewEntryType('login')
      toast.success('Entry created successfully')
    },
    onError: () => {
      toast.error('Failed to create entry')
    },
  })

  const deleteMutation = useMutation({
    mutationFn: (id: string) => apiClient.deleteEntry(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['entries', vaultId] })
      toast.success('Entry deleted successfully')
    },
    onError: () => {
      toast.error('Failed to delete entry')
    },
  })

  const favoriteMutation = useMutation({
    mutationFn: (id: string) => apiClient.toggleFavorite(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['entries', vaultId] })
    },
  })

  const handleCreate = () => {
    if (newEntryTitle.trim()) {
      createMutation.mutate({
        title: newEntryTitle.trim(),
        type: newEntryType,
        data: JSON.stringify({}),
      })
    }
  }

  const copyToClipboard = (text: string) => {
    navigator.clipboard.writeText(text)
    toast.success('Copied to clipboard')
  }

  const filteredEntries = data?.entries.filter((entry) =>
    entry.title.toLowerCase().includes(searchQuery.toLowerCase())
  )

  if (isLoading) {
    return (
      <div className="flex h-full items-center justify-center">
        <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
      </div>
    )
  }

  return (
    <div className="p-6">
      <div className="mb-6">
        <Button
          variant="ghost"
          size="sm"
          className="mb-4"
          onClick={() => navigate('/vaults')}
        >
          <ArrowLeft className="mr-2 h-4 w-4" />
          Back to vaults
        </Button>
        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-3xl font-bold">{vault?.name}</h1>
            <p className="text-muted-foreground">
              {data?.total || 0} entries
            </p>
          </div>
          <Button onClick={() => setIsCreateOpen(true)}>
            <Plus className="mr-2 h-4 w-4" />
            New Entry
          </Button>
        </div>
      </div>

      {data?.entries.length ? (
        <>
          <div className="mb-4">
            <div className="relative">
              <Search className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
              <Input
                placeholder="Search entries..."
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                className="pl-10"
              />
            </div>
          </div>

          <div className="space-y-2">
            {filteredEntries?.map((entry) => (
              <div
                key={entry.id}
                className="group flex items-center gap-4 rounded-lg border bg-card p-4 transition-shadow hover:shadow-md"
              >
                <div className="flex h-10 w-10 items-center justify-center rounded-lg bg-muted">
                  {entryTypeIcons[entry.type] || <FileText className="h-4 w-4" />}
                </div>
                <div className="flex-1 cursor-pointer" onClick={() => copyToClipboard(entry.title)}>
                  <div className="flex items-center gap-2">
                    <h3 className="font-semibold">{entry.title}</h3>
                    {entry.favorite && (
                      <Star className="h-4 w-4 fill-yellow-500 text-yellow-500" />
                    )}
                  </div>
                  <p className="text-sm text-muted-foreground">
                    {entry.type.charAt(0).toUpperCase() + entry.type.slice(1)}
                  </p>
                </div>
                <div className="flex items-center gap-2 opacity-0 transition-opacity group-hover:opacity-100">
                  <Button
                    variant="ghost"
                    size="icon"
                    onClick={() => favoriteMutation.mutate(entry.id)}
                  >
                    <Star
                      className={`h-4 w-4 ${
                        entry.favorite ? 'fill-yellow-500 text-yellow-500' : ''
                      }`}
                    />
                  </Button>
                  <DropdownMenu>
                    <DropdownMenuTrigger asChild>
                      <Button variant="ghost" size="icon">
                        <MoreHorizontal className="h-4 w-4" />
                      </Button>
                    </DropdownMenuTrigger>
                    <DropdownMenuContent align="end">
                      <DropdownMenuItem onClick={() => copyToClipboard(entry.title)}>
                        <Copy className="mr-2 h-4 w-4" />
                        Copy title
                      </DropdownMenuItem>
                      <DropdownMenuSeparator />
                      <DropdownMenuItem
                        className="text-destructive focus:text-destructive"
                        onClick={() => {
                          if (confirm('Are you sure you want to delete this entry?')) {
                            deleteMutation.mutate(entry.id)
                          }
                        }}
                      >
                        <Trash2 className="mr-2 h-4 w-4" />
                        Delete
                      </DropdownMenuItem>
                    </DropdownMenuContent>
                  </DropdownMenu>
                </div>
              </div>
            ))}
          </div>
        </>
      ) : (
        <div className="flex h-[50vh] flex-col items-center justify-center rounded-lg border border-dashed">
          <p className="mb-4 text-lg text-muted-foreground">
            No entries yet
          </p>
          <Button onClick={() => setIsCreateOpen(true)}>
            <Plus className="mr-2 h-4 w-4" />
            Create your first entry
          </Button>
        </div>
      )}

      <Dialog open={isCreateOpen} onOpenChange={setIsCreateOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Create New Entry</DialogTitle>
            <DialogDescription>
              Enter a title for your new entry
            </DialogDescription>
          </DialogHeader>
          <div className="space-y-4 py-4">
            <div className="space-y-2">
              <label className="text-sm font-medium">Entry Type</label>
              <div className="flex gap-2">
                {(['login', 'note', 'card', 'identity'] as const).map((type) => (
                  <Button
                    key={type}
                    variant={newEntryType === type ? 'default' : 'outline'}
                    size="sm"
                    onClick={() => setNewEntryType(type)}
                    className="flex-1"
                  >
                    {entryTypeIcons[type]}
                    <span className="ml-2 capitalize">{type}</span>
                  </Button>
                ))}
              </div>
            </div>
            <div className="space-y-2">
              <label className="text-sm font-medium">Title</label>
              <Input
                placeholder="Entry title"
                value={newEntryTitle}
                onChange={(e) => setNewEntryTitle(e.target.value)}
                onKeyDown={(e) => {
                  if (e.key === 'Enter') handleCreate()
                }}
              />
            </div>
          </div>
          <DialogFooter>
            <Button variant="outline" onClick={() => setIsCreateOpen(false)}>
              Cancel
            </Button>
            <Button
              onClick={handleCreate}
              disabled={!newEntryTitle.trim() || createMutation.isPending}
            >
              {createMutation.isPending && (
                <Loader2 className="mr-2 h-4 w-4 animate-spin" />
              )}
              Create
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  )
}
