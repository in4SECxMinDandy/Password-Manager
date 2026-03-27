import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { apiClient } from '@/lib/api-client'
import { useState } from 'react'
import { Plus, Loader2, MoreHorizontal, Trash2, Edit } from 'lucide-react'
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
import { useNavigate } from 'react-router-dom'

export function VaultsPage() {
  const queryClient = useQueryClient()
  const navigate = useNavigate()
  const [isCreateOpen, setIsCreateOpen] = useState(false)
  const [newVaultName, setNewVaultName] = useState('')

  const { data, isLoading } = useQuery({
    queryKey: ['vaults'],
    queryFn: () => apiClient.getVaults(),
  })

  const createMutation = useMutation({
    mutationFn: (name: string) => apiClient.createVault({ name }),
    onSuccess: (newVault) => {
      queryClient.invalidateQueries({ queryKey: ['vaults'] })
      setIsCreateOpen(false)
      setNewVaultName('')
      toast.success('Vault created successfully')
      navigate(`/vaults/${newVault.id}/entries`)
    },
    onError: () => {
      toast.error('Failed to create vault')
    },
  })

  const deleteMutation = useMutation({
    mutationFn: (id: string) => apiClient.deleteVault(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['vaults'] })
      toast.success('Vault deleted successfully')
    },
    onError: () => {
      toast.error('Failed to delete vault')
    },
  })

  const handleCreate = () => {
    if (newVaultName.trim()) {
      createMutation.mutate(newVaultName.trim())
    }
  }

  if (isLoading) {
    return (
      <div className="flex h-full items-center justify-center">
        <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
      </div>
    )
  }

  return (
    <div className="p-6">
      <div className="mb-6 flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold">My Vaults</h1>
          <p className="text-muted-foreground">
            Manage your password vaults
          </p>
        </div>
        <Button onClick={() => setIsCreateOpen(true)}>
          <Plus className="mr-2 h-4 w-4" />
          New Vault
        </Button>
      </div>

      {!data?.vaults.length ? (
        <div className="flex h-[50vh] flex-col items-center justify-center rounded-lg border border-dashed">
          <p className="mb-4 text-lg text-muted-foreground">
            No vaults yet
          </p>
          <Button onClick={() => setIsCreateOpen(true)}>
            <Plus className="mr-2 h-4 w-4" />
            Create your first vault
          </Button>
        </div>
      ) : (
        <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
          {data.vaults.map((vault) => (
            <div
              key={vault.id}
              className="group relative rounded-lg border bg-card p-4 transition-shadow hover:shadow-md"
            >
              <div className="flex items-start justify-between">
                <div
                  className="flex-1 cursor-pointer"
                  onClick={() => navigate(`/vaults/${vault.id}/entries`)}
                >
                  <h3 className="font-semibold">{vault.name}</h3>
                  <p className="text-sm text-muted-foreground">
                    Created {new Date(vault.created_at).toLocaleDateString()}
                  </p>
                </div>
                <DropdownMenu>
                  <DropdownMenuTrigger asChild>
                    <Button
                      variant="ghost"
                      size="icon"
                      className="h-8 w-8 opacity-0 transition-opacity group-hover:opacity-100"
                    >
                      <MoreHorizontal className="h-4 w-4" />
                    </Button>
                  </DropdownMenuTrigger>
                  <DropdownMenuContent align="end">
                    <DropdownMenuItem
                      onClick={() =>
                        navigate(`/vaults/${vault.id}/entries`)
                      }
                    >
                      <Edit className="mr-2 h-4 w-4" />
                      Open
                    </DropdownMenuItem>
                    <DropdownMenuSeparator />
                    <DropdownMenuItem
                      className="text-destructive focus:text-destructive"
                      onClick={() => {
                        if (
                          confirm('Are you sure you want to delete this vault?')
                        ) {
                          deleteMutation.mutate(vault.id)
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
      )}

      <Dialog open={isCreateOpen} onOpenChange={setIsCreateOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Create New Vault</DialogTitle>
            <DialogDescription>
              Enter a name for your new vault
            </DialogDescription>
          </DialogHeader>
          <div className="py-4">
            <Input
              placeholder="Vault name"
              value={newVaultName}
              onChange={(e) => setNewVaultName(e.target.value)}
              onKeyDown={(e) => {
                if (e.key === 'Enter') handleCreate()
              }}
            />
          </div>
          <DialogFooter>
            <Button variant="outline" onClick={() => setIsCreateOpen(false)}>
              Cancel
            </Button>
            <Button
              onClick={handleCreate}
              disabled={!newVaultName.trim() || createMutation.isPending}
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
